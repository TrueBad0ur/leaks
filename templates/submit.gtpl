<html>
    <head>
    <title></title>
    </head>
    <body>
        <form action="/submit" method="post">
            Your data: <input type="text" name="data">
            <input type="submit" value="submit">
        </form>
        {{ if .outputText }}
            {{ .outputText }}
        {{ end }}
    </body>
</html>