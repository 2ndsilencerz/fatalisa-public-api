# fatalisa-public-api - QRIS Parser API

API provided to parse QRIS used in Indonesia

API list

    /api/qris
            /mpm/{raw_qris} GET
            /mpm            POST (json request { "raw" : "raw_qris" })
            /cpm            POST (json request { "raw" : "raw_qris" })