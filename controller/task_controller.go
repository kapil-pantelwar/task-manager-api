package controller

import (
	"encoding/json"
	//"log"
	"net/http"
	"strconv"
	"task-manager/usecase"

	"github.com/go-chi/chi/v5"
)


type TaskController struct {
	uc *usecase.TaskUseCase
	authUC *usecase.AuthUseCase
}

func NewTaskController(taskUC *usecase.TaskUseCase, authUC *usecase.AuthUseCase) *TaskController{
	return &TaskController{uc:taskUC, authUC: authUC}
}

//Login handles POST /login (fake login for demo)
func (c *TaskController) Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
		return 
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return 
	}

	sessionID, err := c.authUC.Login(input.Username, input.Password)
	if err != nil {
		//log.Println("Error at controller side...") //debug
		http.Error(w,err.Error(), http.StatusUnauthorized)
		return 
	}
	http.SetCookie(w,&http.Cookie{
		Name: "session_id",
		Value: sessionID,
		Path: "/",
		HttpOnly: true,
		MaxAge: 3600, // 1 hour
	})

	json.NewEncoder(w).Encode(map[string]string{"message":"Logged in"})
}

func (c *TaskController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w,"No session found",http.StatusBadRequest)
		return 
	}

	if err := c.authUC.Logout(cookie.Value); err != nil {
		http.Error(w,"Failed to logout", http.StatusInternalServerError)
		return 
	}
	http.SetCookie(w,&http.Cookie{
		Name: "session_id",
		Value: "",
		Path: "/",
		HttpOnly: true,
		MaxAge: -1,
	})
json.NewEncoder(w).Encode(map[string]string{"message":"Logged out"})
}

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title		*string `json:"title"`
		Description *string `json:"description"`
		Status		*string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || (input.Title == nil) || (input.Description == nil) || (input.Status == nil) {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return 
	}

	task, err:= c.uc.Create(input.Title, input.Description,input.Status)
	if err != nil {
		http.Error(w, err.Error(),http.StatusBadRequest)
		return 
	}
	// removed w.Header().Set("Content-Type", "application/json") //sending data in json
	w.WriteHeader(http.StatusCreated) //sets the status to 201
	json.NewEncoder(w).Encode(task) //write into json
}

//GetAlltasks handles GET /tasks

func (c *TaskController) GetAllTasks(w http.ResponseWriter,r *http.Request){
	tasks, err:= c.uc.GetAll()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type:", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//GetTaskByID hanles GET /tasks/{id}

func (c *TaskController) GetTaskByID(w http.ResponseWriter, r *http.Request){
	idStr := chi.URLParam(r,"id")
	id,err:= strconv.Atoi(idStr)
	if err != nil {
		http.Error(w,"Invalid task ID",http.StatusBadRequest)
		return
	}

	task, err := c.uc.GetByID(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return 
	}

	//w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(task)
}

//UpdateTask hanldes PATCH  /tasks/{id}  -- replaced PUT

func (c *TaskController) UpdateTask(w http.ResponseWriter,r *http.Request){
	idStr := chi.URLParam(r,"id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID",http.StatusBadRequest)
	return }

	var input struct {
		Title 		*string `json:"title"`
		Description *string `json:"description"`
		Status 		*string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return 
	}

	task, err := c.uc.Update(id, input.Title,input.Description,input.Status)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return 
	}

	//w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(task)

}

//DeleteTask handles DELETE /tasks/{id}

func (c *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r,"id")
	id,err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w,"Invalid Task ID",http.StatusBadRequest)
		return 
	}

	cookie, err:= r.Cookie("session_id")
	if err != nil {
		http.Error(w,"Unauthorized",http.StatusUnauthorized)
		return
	}

	authorized,err := c.authUC.Authorize(cookie.Value,"admin")
	if err != nil || !authorized {
		http.Error(w,"Forbidden",http.StatusForbidden)
		return
	}

	if err := c.uc.Delete(id); err != nil {
		http.Error(w,err.Error(),http.StatusNotFound)
		return 
	}
	w.WriteHeader(http.StatusOK) //successfully deleted the task by authorized user
}


