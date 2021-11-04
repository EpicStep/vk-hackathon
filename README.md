# VK Hackathon (Back End)
## Run (Locally)
Create .env file and add following values:
```dotenv
PORT=8181
DATABASE_URL=root:mypass@tcp(mysql:3306)/test
```

After you can run app in docker:
```bash
docker compose up -d
```
You need to insert ```migrations/000001_initial.up.sql```

Also, you can run tests:
```bash
bash test.bash
```