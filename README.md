# Local NAS

simple NAS web app for personal usage with authentication and authorization

### Environment Variables

```
############Mysql#################
DBPWD     -> Database Password    
DBUSR     -> Database Username    
DBADDR    -> Database Address     
DBNAME    -> Database Name        
############Redis#################
REDISPWD  -> redis password       
REDISADDR -> redis address        
############General###############
PORT      -> Listen Port          
##################################
```

### Constant Variable

- ```rootdir``` -> Root Directory for Uploaded Files
 - /upload/upload.go at line 16


#### Database

- Credentials : **Mysql**
- Session : **Redis**

**tables** :

```sql
CREATE TABLE cloud (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(40) UNIQUE NOT NULL,
  password VARCHAR(32) NOT NULL
);

INSERT INTO cloud(username, password) VALUES('your username', MD5('your password'));
```
