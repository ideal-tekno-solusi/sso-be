# sso-be

TODO:
- [x] create endpoint authorization
- [x] create logic to check user already login by check Authorization cookie
- [ ] create endpoint login
- [ ] create logic to decrypt encrypted pass
- [ ] create logic to hash pass (bcrypt)
- [ ] create logic to check csrf from cookie and req is match
- [ ] create logic to generate new csrf fo login form
- [ ] create logic to veridy code challenge oauth 2.0

# note
after edit query in database/postgresql/query.sql, dont forget to run
`sqlc generate`