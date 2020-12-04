package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)
type Server struct {
	subjects map[string]map[string]float64
	students map[string]map[string]float64
}

type Data struct {
	Student string
	Subject string
	Grade   float64
}

var server Server

func (this *Server) AddGrade(data Data) bool {

	if _, ok := this.students[data.Student]; !ok {
		grade := make(map[string]float64)
		grade[data.Subject] = data.Grade
		this.students[data.Student] = grade
	} else {
		if _, ok := this.students[data.Student][data.Subject]; ok {
			return false
		}
		this.students[data.Student][data.Subject] = data.Grade
	}

	if _, ok := this.subjects[data.Subject]; !ok {
		grade := make(map[string]float64)
		grade[data.Student] = data.Grade
		this.subjects[data.Subject] = grade
	} else {
		this.subjects[data.Subject][data.Student] = data.Grade
	}
	return true
}
func (this *Server) SubjectAVG(data Data) float64{
	avg := float64(0)
	for student := range this.subjects[data.Subject] {
		avg += this.subjects[data.Subject][student]
	}
	avg = avg / float64(len(this.subjects[data.Subject]))

	return avg
}

func (this *Server) StudentAVG(data Data) float64 {
	avg := GetStudentAVG(data.Student, this)
	return avg
}

func (this *Server) GeneralAVG() float64 {

	avg := float64(0)
	for student := range this.students {
		avg += GetStudentAVG(student, this)
	}
	avg = avg / float64(len(this.students))
	return avg
}

func GetStudentAVG(student string, this *Server) float64 {
	avg := float64(0)
	for subject := range this.students[student] {
		avg += this.students[student][subject]
	}
	avg = avg / float64(len(this.students[student]))
	return avg
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("index.html"),
	)
}

func addGrade(res http.ResponseWriter, req *http.Request) {

	fullHTML := cargarHtml("./addGrade.html")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		form := req.PostForm
		if form["grade"][0] != "" && form["student"][0] != "" && form["subject"][0] != "" {
			var data Data
			data.Grade, _ = strconv.ParseFloat(form["grade"][0], 64)
			data.Student = form["student"][0]
			data.Subject = form["subject"][0]

			status := server.AddGrade(data)
			if status {
				fullHTML += "<h4>Calificacion agregada</h4></body></html>"
			} else {
				fullHTML += "<h4>El alumno ya cuenta con calificacion para la materia</h4></body></html>"
			}
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		} else {
			fullHTML += "</body></html>"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		}
	case "GET":
		fullHTML += "</body></html>"
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			fullHTML,
		)
	}
}
func studentAVG(res http.ResponseWriter, req *http.Request) {
	
	fullHTML := cargarHtml("./AVGStudent.html")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		form := req.PostForm
		if form["student"][0] != "" {
			var data Data
			data.Student = form["student"][0]
			avg := server.StudentAVG(data)

			fullHTML = fullHTML + "<h4>Estudiante: " + data.Student + "</h4>" +
			 "<h4>Promedio: " + strconv.FormatFloat(avg, 'f', -1,64) + "</h4>" +
			 "</body> </html>"

			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		} else {
			fullHTML += "</body></html>"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		}
	case "GET":
			fullHTML += "</body></html>"
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			fullHTML,
		)
	}
}

func subjectAVG(res http.ResponseWriter, req *http.Request) {
	
	fullHTML := cargarHtml("./AVGSubject.html")
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		form := req.PostForm
		if form["subject"][0] != "" {
			var data Data
			data.Subject = form["subject"][0]
			avg := server.SubjectAVG(data)

			fullHTML = fullHTML + "<h4>Materia: " + data.Subject + "</h4>" +
			 "<h4>Promedio: " + strconv.FormatFloat(avg, 'f', -1,64) + "</h4>" +
			 "</body> </html>"

			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		} else {
			fullHTML += "</body></html>"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		}
	case "GET":
			fullHTML += "</body></html>"
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			fullHTML,
		)
	}
}

func generalAVG(res http.ResponseWriter, req *http.Request) {
	
	fullHTML := cargarHtml("./AVGGeneral.html")
	switch req.Method {
	case "GET":
			avg := server.GeneralAVG()
			fullHTML = fullHTML + strconv.FormatFloat(avg, 'f', -1,64) + "</h4>" +
			 "</body> </html>"

			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				fullHTML,
			)
		}
}

func main() {
	server.subjects = make(map[string]map[string]float64)
	server.students = make(map[string]map[string]float64)
	http.HandleFunc("/", root)
	http.HandleFunc("/addGrade", addGrade)
	http.HandleFunc("/studentAVG", studentAVG)
	http.HandleFunc("/subjectAVG", subjectAVG)
	http.HandleFunc("/generalAVG", generalAVG)
	fmt.Println("Arrancando el servidor...")
	http.ListenAndServe(":9000", nil)
}
