
# BC Project
Hi! This project created for combining an analysis application with API.


##  Built With
* [Golang](https://go.dev/)
* [Amazon Web Services](https://aws.amazon.com/tr/)
* [Docker](https://www.docker.com/)
* [PostgreSQL](postgresql.org)
* [MongoDB](https://www.mongodb.com/)
* [JSON Web Tokens](https://jwt.io/)
* [R-Shiny](https://shiny.rstudio.com/)
* [Bootstrap](https://getbootstrap.com)
* [HTML](https://html.com/)


## Usage

### [Login Page](http://18.206.88.12:9000/login)

<img src="https://cdn.discordapp.com/attachments/519918508998656028/929023335118041138/unknown.png" alt="Logo" width="50%">

The Dockerized application is served on AWS EC2 server you can test login page with username : "testCafe" and password : "1234".
When you submit the username and password, Golang connect to MongoDB and search if is there any record with this username. If there is a match, Golang API will hash the submitted password and compare with password that came from MongoDB. You can try with test file which in temp file and file name is "client data.csv"

### [Welcome Page](http://18.206.88.12:9000/welcome)

<img src="https://media.discordapp.net/attachments/519918508998656028/929025950585352302/unknown.png" alt="Logo" width="50%">

After the login succesfully, the Golang API gives you an authentication token and redirect you to Welcome Page. If user try to pass the login page and connect directly to Welcome Page. Because of missing token you can't connect this page. 

In this page you can upload formatted .csv files for sending to analysis system.

#### Sample of formatted .csv file
| orderID| sentBy| orderTime| itemVal| closedBy| restName| paidBy| from| item| itemCount| table| closedAt |
 |-----------|:-----------:|-----------:| :-----------:|:-----------:|:-----------:|:-----------:|:-----------:|:-----------:|:-----------:|:-----------:|:-----------:|
| 1 | P1 | 2021-09-21 15:22:07 | 29.67 | testCafeAdmin | testCafe | cash|admin | CheeseCake | 1 | d14 | 2021-09-21 15:59:31| 
|2|P1|2021-08-13 15:52:51|17.82|testCafeAdmin|testCafe|card|admin|Limonata|1|d10|2021-08-13 16:32:58|


If you upload correct typed and formatted files, you will see a confirmation notification and you will see a new button which will redirect you to R-Shiny Page. This page also worked on AWS EC2 services.

<img src="https://cdn.discordapp.com/attachments/519918508998656028/929030537287434310/unknown.png" alt="Logo" width="40%">

### [R-Shiny Page](http://3.145.16.200:3838/project_bc)

In the end, both data which came from database and which provided by clients combined and send to the analysis service. Both data visualized together.

<img src="https://media.discordapp.net/attachments/519918508998656028/929031871470399538/unknown.png?width=732&height=671" alt="Logo" width="40%">
