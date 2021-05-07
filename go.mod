module github.com/polarspetroll/localcloud

go 1.15

require (
  login v0.0.0
  upload v0.0.0
  )

replace upload => ./upload
replace login => ./login
