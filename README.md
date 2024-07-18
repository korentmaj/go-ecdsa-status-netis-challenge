# netis.si_a9af8eb9dc9831b00dafec4a2f381dd98d18b5e12c65e454cc24fa9a4c2cb0c7
Pričakovan programski jezik: Go

1. 
Napiši sledeče funkcije:
    a.) Funkcija, ki kreira ecdsa-p256 ključ in ga zapiše v PEM obliki.	
    b.) Funkcija, ki prebere PEM ključ in z njim podpiše sporočilo. Podpis vrne v Base64URL obliki. 
    c.) Funkcija, ki prebere PEM ključ in z njim preveri podpis sporočila iz druge (b.) funkcije.
	
2. 
Napiši strukturo za hranjenje enega ali več statusov (boolean true/false). 
Struktura naj uporabi tip []byte za hranjenje statusov, kjer vsak bit hrani eno stanje (true ali false). 
Struktura naj podpira nastavljanje posameznega obstoječega stanja, kreacijo novega stanja v seznamu ter enkodiranje celotnega seznama v gzipped base64 enkodiranem nizu; encodedList: (return base64(gzip(byteArray))).

3. 
Sprogramiraj HTTP strežnik, ki servira JSON REST API-je, preko katerih lahko kreiraš in manipuliraš s strukturo iz 2. naloge. Vsak seznam naj se shrani v PostgreSQL bazo. Seznam API klicev naj bo sledeč:

------------------------- API -----------------------------
    GET /api/status/:statusId#{{index}}
		- Spremenljivka index je tipa integer.
        - Vrne JWS Compact podpisano sporočilo, z JWK headerjem in sledečim primerom payload-a:
            {
                "iat": {{UNIX_time_now}}
                "exp": {{UNIX_time_now+1day}}
                "iss": "http://{{domain}}/api/status/{{:statusId}}"
                "status": {
                    "encodedList": {{encodedList}}
                    "index": {{index}}
                }
            }
		- Pravilnost JWS-ja (podpisanega JWT-ja), lahko preveriš tukaj: https://jwt.io/ 
    PUT /api/status/:statusId#{{index}}
        - Nastavi stanje {{index}} v strukturi {{:statusId}} na true.
    DELETE /api/status/:statusId#index
        - Nastavi stanje {{index}} v strukturi {{:statusId}} na false.
    POST /api/status/:statusId
        - Kreiraj novo stanje v strukturi {{:statusId}} in vrni index.
        
    GET /api/status
        - Dobi vse strukture (seznam status Id-jev)
    POST /api/status
        - Kreiraj novo strukturo in vrni Id nove strukture.
------------------------ /API ------------------------------

BONUS:
    a.) Napiši funkcijo, ki kot argument prejme URL (to je HTTP GET API, npr: http://localhost:8000/api/status/a_Q5JxCz#1), izvede klic na ta API in prejme JWS - podpisan JWT.
		Nato preveri podpis prejetega JWS-ja (preveri lahko tudi iat, exp in iss headerje), prebere status iz payload-a in na koncu vrne boolean stanje (true ali false) index-a v encodedList-i.
    
    b.) Nadgradi 3.nalogo, da dodaš avtorizacijo na vse PUT, POST in DELETE http metode (lahko je osnovna preko Authorization header-ja ali pa kaj bolj naprednega).
    
    c.) Enkriptiraj sezname stanj v bazi (način enkriptiranja je lahko AES password, ali pa s kreacijo RSA ključa).
	
	
-------------------------------------------------------------

Rešitve nalog se pošljejo kot .zip Go modul projekta, ki naj vsebuje /database mapo s potrebnimi SQL stavki za rekreiranje rešitev. 
