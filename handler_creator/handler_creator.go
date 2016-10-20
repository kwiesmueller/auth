package handler_creator

import (
	"fmt"
	"net/http"

	"github.com/bborbe/auth/directory/application_directory"
	"github.com/bborbe/auth/directory/group_user_directory"
	"github.com/bborbe/auth/directory/token_user_directory"
	"github.com/bborbe/auth/directory/user_data_directory"
	"github.com/bborbe/auth/directory/user_group_directory"
	"github.com/bborbe/auth/directory/user_token_directory"
	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/router"
	"github.com/bborbe/auth/service/application"
	"github.com/bborbe/auth/service/user"
	"github.com/bborbe/auth/service/user_data"
	"github.com/bborbe/auth/service/user_group"
	"github.com/bborbe/auth/v1"
	"github.com/bborbe/auth/v1/handler_creator"
	"github.com/bborbe/http_handler/check"
	"github.com/bborbe/http_handler/not_found"
	"github.com/bborbe/password/generator"
	"github.com/bborbe/redis_client/ledis"
	"github.com/golang/glog"
)

type HandlerCreator interface {
	CreateHandler(
		prefix string,
		authApplicationPassword string,
		ledisdbAddress string,
		ledisdbPassword string,
	) (http.Handler, error)
}

type handlerCreator struct {
}

func New() *handlerCreator {
	h := new(handlerCreator)
	return h
}

func (h *handlerCreator) CreateHandler(
	prefix string,
	authApplicationPassword string,
	ledisdbAddress string,
	ledisdbPassword string,
) (http.Handler, error) {
	ledisClient := ledis.New(ledisdbAddress, ledisdbPassword)
	passwordGenerator := generator.New()

	tokenUserDirectory := token_user_directory.New(ledisClient)
	userTokenDirectory := user_token_directory.New(ledisClient)
	applicationDirectory := application_directory.New(ledisClient)
	groupUserDirectory := group_user_directory.New(ledisClient)
	userGroupDirectory := user_group_directory.New(ledisClient)
	userDataDirectory := user_data_directory.New(ledisClient)

	applicationService := application.New(passwordGenerator.GeneratePassword, applicationDirectory)
	userService := user.New(userTokenDirectory, userGroupDirectory, tokenUserDirectory, userDataDirectory)
	userGroupService := user_group.New(userGroupDirectory, groupUserDirectory)
	userDataService := user_data.New(userDataDirectory)

	go func() {
		if _, err := applicationService.CreateApplicationWithPassword(model.AUTH_APPLICATION_NAME, model.ApplicationPassword(authApplicationPassword)); err != nil {
			glog.Warningf("create auth application failed: %v", err)
			return
		}
	}()

	checkHandler := check.New(ledisClient.Ping)
	notFoundHandler := not_found.New()
	v1HandlerCreator := handler_creator.New()
	v1Handler, err := v1HandlerCreator.CreateHandler(
		fmt.Sprintf("%s/api/%s", prefix, v1.VERSION),
		applicationService,
		userService,
		userGroupService,
		userDataService,
	)
	if err != nil {
		return nil, err
	}

	handler := router.New(
		prefix,
		notFoundHandler.ServeHTTP,
		checkHandler.ServeHTTP,
		v1Handler.ServeHTTP,
	)
	return handler, nil
}
