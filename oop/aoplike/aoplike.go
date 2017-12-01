package main

import (
	"errors"
	"fmt"
	"strings"
)

//
// Generic aspect wrapper definition
//

type joinPointCompleter func(e error, results ...interface{})

type joinPointCreator func(methodName string, arguments ...interface{}) joinPointCompleter

//
// Sample service interface
//

type userProfile struct {
	id   string
	name string
	age  int
}

type userService interface {
	add(p *userProfile) (string, error)
	getByID(id string) (*userProfile, error)
}

//
// Sample service wrapper (can be generated)
//

type userServiceWrapper struct {
	userService

	jpCreator joinPointCreator
	proxy     userService
}

func (w *userServiceWrapper) add(p *userProfile) (string, error) {
	jpCompleter := w.jpCreator("UserService.add", p)
	id, error := w.proxy.add(p)
	jpCompleter(error, id)
	return id, error
}

func (w *userServiceWrapper) getByID(id string) (*userProfile, error) {
	jpCompleter := w.jpCreator("UserService.getByID", id)
	userProfile, error := w.proxy.getByID(id)
	jpCompleter(error, userProfile)
	return userProfile, error
}

func newUserServiceWrapper(p userService, j joinPointCreator) userService {
	return &userServiceWrapper{jpCreator: j, proxy: p}
}

//
// Sample service itself
//

type defaultUserService struct {
	userService

	users []userProfile
}

func (t *defaultUserService) add(p *userProfile) (string, error) {
	id := p.id

	return id, nil
}

func (t *defaultUserService) getByID(id string) (*userProfile, error) {
	for _, u := range t.users {
		if strings.Compare(id, u.id) == 0 {
			return &u, nil
		}
	}

	return nil, errors.New("No such user")
}

func newDefaultUserService() userService {
	return &defaultUserService{users: []userProfile{}}
}

//
// Usage sample
//

func main() {
	fmt.Println("AOP-Like: Demo of cross-cutting concerns in Golang")

	svc := newUserServiceWrapper(newDefaultUserService(), func(m string, args ...interface{}) joinPointCompleter {
		fmt.Printf("[Aspect] Before method: %s\n\targs=%s\n", m, args)
		return func(e error, results ...interface{}) {
			fmt.Printf("[Aspect] After method:  %s\n\te=%s args=%s\n", m, e, results)
		}
	})

	id, err := svc.add(&userProfile{name: "bob", id: "1", age: 15})
	if err != nil {
		fmt.Printf("Error while adding bob, e=%s\n", err)
	} else {
		fmt.Printf("Added user bob, id=%s\n", id)
	}

	p, err := svc.getByID(id)
	if err != nil {
		fmt.Printf("Unable to locate user bob, e=%s\n", err)
	} else {
		fmt.Printf("Found user bob: {name: '%s', age: %d}\n", p.name, p.age)
	}

	fmt.Println("---")
}
