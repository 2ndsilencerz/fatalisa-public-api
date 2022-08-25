# fatalisa-public-api - Pray Schedule API

API provided to serve pray schedule based on city and date

API List

    /api
        /pray-schedule           POST (json request { "city" : "city", "date" : "YYYY/MM/DD" })
        /pray-schedule/city-list GET
        /pray-schedule/{city}    GET