package router

import (
	"net/http"

	"fmt"

	"github.com/bborbe/http_handler/adapter"
	"github.com/gorilla/mux"
)

type handler struct {
	router http.Handler
}

func New(
	prefix string,
	notFound http.HandlerFunc,
	version http.HandlerFunc,
	login http.HandlerFunc,
	applicationCreate http.HandlerFunc,
	applicationDelete http.HandlerFunc,
	applicationGet http.HandlerFunc,
	userRegister http.HandlerFunc,
	userUnregister http.HandlerFunc,
	userDelete http.HandlerFunc,
	tokenAdd http.HandlerFunc,
	tokenRemove http.HandlerFunc,
	userGroupAdd http.HandlerFunc,
	userGroupRemove http.HandlerFunc,
	userDataSet http.HandlerFunc,
	userDataSetValue http.HandlerFunc,
	userDataGet http.HandlerFunc,
	userDataGetValue http.HandlerFunc,
	userDataDelete http.HandlerFunc,
	userDataDeleteValue http.HandlerFunc,
	userList http.HandlerFunc,
) *handler {
	router := mux.NewRouter()

	router.Path(fmt.Sprintf("%s/version", prefix)).Methods("GET").HandlerFunc(version)

	router.Path(fmt.Sprintf("%s/user", prefix)).Methods("GET").HandlerFunc(userList)

	router.Path(fmt.Sprintf("%s/user/{username}/data", prefix)).Methods("POST").HandlerFunc(userDataSet)
	router.Path(fmt.Sprintf("%s/user/{username}/data/{key}", prefix)).Methods("POST").HandlerFunc(userDataSetValue)
	router.Path(fmt.Sprintf("%s/user/{username}/data", prefix)).Methods("GET").HandlerFunc(userDataGet)
	router.Path(fmt.Sprintf("%s/user/{username}/data/{key}", prefix)).Methods("GET").HandlerFunc(userDataGetValue)
	router.Path(fmt.Sprintf("%s/user/{username}/data", prefix)).Methods("DELETE").HandlerFunc(userDataDelete)
	router.Path(fmt.Sprintf("%s/user/{username}/data/{key}", prefix)).Methods("DELETE").HandlerFunc(userDataDeleteValue)

	router.Path(fmt.Sprintf("%s/login", prefix)).Methods("POST").HandlerFunc(login)

	router.Path(fmt.Sprintf("%s/application", prefix)).Methods("POST").HandlerFunc(applicationCreate)
	router.PathPrefix(fmt.Sprintf("%s/application/", prefix)).Methods("DELETE").HandlerFunc(applicationDelete)
	router.PathPrefix(fmt.Sprintf("%s/application/", prefix)).Methods("GET").HandlerFunc(applicationGet)

	router.Path(fmt.Sprintf("%s/user/{username}", prefix)).Methods("DELETE").HandlerFunc(userDelete)
	router.Path(fmt.Sprintf("%s/user", prefix)).Methods("POST").HandlerFunc(userRegister)

	router.Path(fmt.Sprintf("%s/token/{token}", prefix)).Methods("DELETE").HandlerFunc(userUnregister)
	router.Path(fmt.Sprintf("%s/token", prefix)).Methods("POST").HandlerFunc(tokenAdd)
	router.Path(fmt.Sprintf("%s/token", prefix)).Methods("DELETE").HandlerFunc(tokenRemove)

	router.Path(fmt.Sprintf("%s/user_group", prefix)).Methods("POST").HandlerFunc(userGroupAdd)
	router.Path(fmt.Sprintf("%s/user_group", prefix)).Methods("DELETE").HandlerFunc(userGroupRemove)

	router.NotFoundHandler = adapter.New(notFound)

	h := new(handler)
	h.router = router
	return h
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(resp, req)
}
