package email

import "fmt"

func Register() string {
	return fmt.Sprint(`
	<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Cabin:wght@400;500&display=swap" rel="stylesheet">
</head>
<body>
    <div style="
        width: 50%;
        margin: 30px auto 0 auto;">
        <h1 style="
            font-size: 1.5rem;
            font-weight: 500;
            margin: 0 auto;
            font-family: 'Arial';
            ">
            Hai, {{.name}}!
        </h1><br/>
        <p style="
            margin: 0 auto;
            font-size: 1.2rem;
            font-weight: 400;
            font-family: 'Arial';">
            Kamu baru saja melakukan registrasi di AIVue. Klik tombol di bawah untuk memverifikasi email kamu!
        </p><br/><br/>
        <a href="{{.link}}" 
            target="_blank" 
            style="
            margin: 0 auto;
            padding: 10px;
            display: block;
            width: 180px;
            height: 20px;
            background-color: #7E95FF;
            color: white;
            text-align: center;
            border-radius: 5px;
            text-decoration: none;
            font-family: 'Arial';">
            Verifikasi Email
        </a>
    </div>
</body>
</html>`)
}

func VerifySuccess() string {
	return fmt.Sprint(`
    <!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Cabin:wght@400;500&display=swap" rel="stylesheet">
</head>
<body>
    <div style="
        width: 50%;
        margin: 30px auto 0 auto;">
        <h1 style="
            font-size: 1.5rem;
            font-weight: 500;
            margin: 0 auto;
            font-family: 'Arial';
            ">
            Hai, {{.name}}!
        </h1><br/>
        <p style="
            margin: 0 auto;
            font-size: 1.2rem;
            font-weight: 400;
            font-family: 'Arial';">
            Selamat datang di AIVue! Berikut adalah password untuk akun Anda: <strong>{{.password}}</strong>
        </p><br/>
    </div>
</body>
</html>`)
}
