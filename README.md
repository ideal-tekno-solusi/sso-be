# sso-be

TODO:
- [x] create endpoint authorization
- [ ] create logic to check user already login by check Authorization cookie
- [ ] create logic to check csrf from cookie and req is match
- [ ] create logic to generate new csrf fo login form
- [ ] create logic to veridy code challenge oauth 2.0

# note
after edit query in database/postgresql/query.sql, dont forget to run
`sqlc generate`