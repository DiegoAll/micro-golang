
INVESTIGAR DE FORMA ADECUADA PARA CREAR MODELOS EN GOLANG
TREVOR PROPONE MODELOS PARA POSTGRES Y ESE TEMA DEL POOL (no entiendo)


ALTERNATIVA INTERESANTE:
MICROSERVICIO DE AUTENTICACION CON DB PROPIA E INDEPENDIENTE (INVENTAR PROTECCION)
LA APLICACION CON EL REPOSITORIO TIPO ESCOTO REPOSITORY + (POSTGRES.GO | MYSQL.GO)
DUMMY DE MUESTRA CON EL CRUD Y EL USUARIO YA SIGUE SU APP. INTENTAR CREAR EL BUSCAR. (OPTIONAL)


HAY N MIL OPORTUNIDADES DE AGREGAR SEGURIDAD.
COMO ES MAS ADECUADO PRESENTARLAS, TODAS SON DIRECTAMENTE EN EL CODIGO Y NO NECESARIAMENTE SON DEL OWASP API SECURITY TOP 10... UUID Y TODO ESO. QUIZAS COMO UNA GUIA DE BUENAS PRACTICAS O SERIA EXPLICAR COMO SE DISEÑO LA APP Y QUE SE IMPLEMENTO PARA APLICAR SEGURIDAD ?


TEMA DE LAS DEPENDENCIAS DIRECTAS, INDIRECTAS Y TRANSITIVAS

front-end
go run ./cmd/web/


broker
go run ./cmd/api/


CURRENT_UID=$(id -u):$(id -g) docker-compose up -d

Vulnerabilidades en la imagen:

https://hub.docker.com/layers/library/golang/1.18beta2-alpine/images/sha256-0538f1f40c327b7612384b746ee9bd4c1ce7c5b58b6191a40605910a904fb5c3


**dannysalcedo servilleta frontend no hace parte de los servicios, o chats teams (nunca contenerizan el front escoto no (validar) trevor no ?)** 
**No se "contenerizan" las DB, depende**

Hay una diferencia al ejecutar a mano la app y con Makefile.

Makefile genera una imagen del servicio a través de Go build.


chown: changing ownership of '/var/lib/postgresql/data': Operation not permitted

ya aplique el chmod 777 tanto para db-data como para db-data/postgres
agregue    user: ${CURRENT_UID} en cada servicio del docker-compose

de govcarpeta  chmod -R 777 db-data


de internet

docker-compose run db ls -la /var/lib/postgresql

```
(base) CO0C02GD0T7MD6M:project dposada$ docker-compose run project-postgres-1 ls -la /var/lib/postgresql
WARN[0000] The "CURRENT_UID" variable is not set. Defaulting to a blank string. 
WARN[0000] The "CURRENT_UID" variable is not set. Defaulting to a blank string.
```

Intente agregando estos parametros al docker-compose y no funciona se cae sigue el tema del permiso

      #CURRENT_UID: (id -u):$(id -g)
      CURRENT_UID: 164865804:1010544492

NO FUE TAN MAMEY COMO EN GOVCARPETA.
LE DI PERMISOS A TODA LA CARPETA PROJECT.

FUNCIONO AL PARECER CON EN EL MAKEFILE  	
CURRENT_UID=164865804:1010544492 docker-compose up --build -d

AL CAMBIAR docker-compose up --build -d VOLVIO A APARECER EL CAMPO PARA INGRESAR LA CONTRASEÑA, con docker-compose up -d (Creo que es normal)

Siempre habia sido un tema de permisos tipo chmod

INDAGAR BIEN QUE GONORREA DE ORDINARIES QUEMAR ESTO O PROBAR EN ALGO CON LINUX PARA CURARSE EN SALUD PARA LA TESIS.

AL HABER QUEMADO ESTE USUARIO SI HUBO AL PARECER PERMISOS Y SE LLENO LA CARPETA DE POSTGRES.

quite este     user: ${CURRENT_UID}  del docker-compose y se volvio a dañar por el permiso

Entonces se puede inferir que cada servicio necesita su propio?

  postgres:
    user: ${CURRENT_UID}


Solo se que se quitaron de las variables de entorno.

SEGUIR PROBANDO...


por que no coge con variables?

por que si funciona para docker-compose down?

para up -d o build no?

