-- mysql :

CREATE TABLE cloud (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(40) UNIQUE NOT NULL,
  password VARCHAR(32) NOT NULL
);
