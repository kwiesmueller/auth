package router

import (
	"net/http"

	"fmt"

	"github.com/bborbe/auth/model"
	"github.com/bborbe/auth/v1"
	"github.com/gorilla/mux"
)

type HandlerCreator interface {
	Prefix() model.Prefix
	NotFoundHandler() http.Handler
	HealthzHandler() http.Handler
	ReadinessHandler() http.Handler
	VersionHandler() http.Handler
	UserListHandler() http.Handler
	UserRegisterHandler() http.Handler
	UserDeleteHandler() http.Handler
	UserDataSetHandler() http.Handler
	UserDataSetValueHandler() http.Handler
	UserDataGetHandler() http.Handler
	UserDataGetValueHandler() http.Handler
	UserDataDeleteHandler() http.Handler
	UserDataDeleteValueHandler() http.Handler
	LoginHandler() http.Handler
	ApplicationCreateHandler() http.Handler
	ApplicationDeleteHandler() http.Handler
	ApplicationGetHandler() http.Handler
	UserUnregisterHandler() http.Handler
	TokenAddHandler() http.Handler
	TokenRemoveHandler() http.Handler
	UserGroupAddHandler() http.Handler
	UserGroupRemoveHandler() http.Handler
	TokensForUsernameHandler() http.Handler
	GroupNamesForUsernameHandler() http.Handler
}

func Create(h HandlerCreator) http.Handler {
	router := mux.NewRouter()

	router.Path(fmt.Sprintf("%s/api/%s/version", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Handler(h.VersionHandler())

	router.Path(fmt.Sprintf("%s/api/%s/user", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Handler(h.UserListHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.UserRegisterHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.UserDeleteHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.UserDataSetHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data/{key}", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.UserDataSetValueHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Handler(h.UserDataGetHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data/{key}", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Handler(h.UserDataGetValueHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.UserDataDeleteHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user/{username}/data/{key}", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.UserDataDeleteValueHandler())

	router.Path(fmt.Sprintf("%s/api/%s/login", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.LoginHandler())

	router.Path(fmt.Sprintf("%s/api/%s/application", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.ApplicationCreateHandler())
	router.PathPrefix(fmt.Sprintf("%s/api/%s/application/", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.ApplicationDeleteHandler())
	router.PathPrefix(fmt.Sprintf("%s/api/%s/application/", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Handler(h.ApplicationGetHandler())

	router.Path(fmt.Sprintf("%s/api/%s/token/{token}", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.UserUnregisterHandler())
	router.Path(fmt.Sprintf("%s/api/%s/token", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Queries("username", "{username}").Handler(h.TokensForUsernameHandler())
	router.Path(fmt.Sprintf("%s/api/%s/token", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.TokenAddHandler())
	router.Path(fmt.Sprintf("%s/api/%s/token", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.TokenRemoveHandler())

	router.Path(fmt.Sprintf("%s/api/%s/user_group", h.Prefix(), v1.VERSION)).Methods(http.MethodGet).Queries("username", "{username}").Handler(h.GroupNamesForUsernameHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user_group", h.Prefix(), v1.VERSION)).Methods(http.MethodPost).Handler(h.UserGroupAddHandler())
	router.Path(fmt.Sprintf("%s/api/%s/user_group", h.Prefix(), v1.VERSION)).Methods(http.MethodDelete).Handler(h.UserGroupRemoveHandler())

	router.Path(fmt.Sprintf("%s/healthz", h.Prefix())).Methods(http.MethodGet).Handler(h.HealthzHandler())
	router.Path(fmt.Sprintf("%s/readiness", h.Prefix())).Methods(http.MethodGet).Handler(h.ReadinessHandler())

	router.NotFoundHandler = h.NotFoundHandler()

	return router
}
