package handlers

import (
	"WIND/internal/userService"
	"WIND/internal/web/users"
	"context"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	Service *userService.UserService
}

func NewUserHandler(service *userService.UserService) *userHandler {
	return &userHandler{
		Service: service,
	}
}

func (u userHandler) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := u.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := users.GetUsers200JSONResponse{}

	for _, tsk := range allUsers {
		user := users.User{
			Id:       &tsk.ID,
			Email:    &tsk.Email,
			Password: &tsk.Password,
		}
		response = append(response, user)
	}
	return response, nil
}

func (u userHandler) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	userRequest := request.Body

	userToCreate := userService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
	}
	createdUser, err := u.Service.CreateUser(userToCreate)
	if err != nil {
		return nil, err
	}

	response := users.PostUsers201JSONResponse{
		Id:       &createdUser.ID,
		Email:    &createdUser.Email,
		Password: &createdUser.Password,
	}
	return response, nil
}

func (u userHandler) DeleteUsersId(_ context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	userId := request.Id

	err := u.Service.DeleteUserById(userId)
	if err != nil {
		if err == userService.ErrUserNotFound {
			return nil, echo.NewHTTPError(404, "User not found")
		}
		return nil, err
	}
	return nil, echo.NewHTTPError(204, "User deleted")
}

func (u userHandler) PatchUsersId(_ context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	userRequest := request.Body

	userToUpdate := userService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
	}

	updatedUser, err := u.Service.UpdateUserById(request.Id, userToUpdate)
	if err != nil {
		return nil, err
	}

	response := users.PatchUsersId200JSONResponse{
		Id:       &updatedUser.ID,
		Email:    &updatedUser.Email,
		Password: &updatedUser.Password,
	}
	return response, nil
}
