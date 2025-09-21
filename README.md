# sso-be

TODO:
- [x] create endpoint authorization
- [x] create logic to check user already login by check Authorization cookie
- [x] update logic to send custom error message when csrf error relate
- [x] create endpoint login
- [x] create endpoint token to exchange auth code
- [x] create logic to decrypt encrypted pass
- [x] create logic to hash pass (bcrypt)
- [x] create logic to check csrf from cookie and req is match
- [x] create logic to generate new csrf fo login form
- [x] create logic to verify code challenge oauth 2.0
- [x] remove csrf and create new api "omni sso" to handle login and GET refresh token, ofcourse with csrf turn on

new TODO:
- [x] rework /login to just validate username and password, then return valid or not, later auth fe that will do request to /auth
- [x] rework /auth
- [x] rework /token to serve token and refresh token
- [x] change all logrus from controller to use one from utils
- [ ] search a method to clean up auth code for refresh if it's not used for too long (7 days)
- [ ] clean up

# note
- after edit query in database/postgresql/query.sql, dont forget to run `sqlc generate`
- this project is made without minding it's securities, this project solely for POC of how fully build enterprise software works internally
- if not in same network, run `cloudflared access tcp --hostname dbstaging.idtecsi.my.id --url localhost:5432` in your terminal host after auth to zero trust