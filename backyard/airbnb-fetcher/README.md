Airbnb-fetcher
==============

Airbn-fetcher es un script de python que recibe como argumento un room id correspondiente a una propiedad en airbnb, y devuelve como output un JSON con el título de la publicación, la descripción en su idioma original, y una lista de imágenes asociadas.

Por ejemplo:

```
# Ejecutando con el room id 34228829
python main.py --room_id 34228829
```

Output:

```
{
  "title": "Cabaña el Ocaso con hermosa vista al mar.",
  "description": "Cabaña El Ocaso ofrece alojamiento en Laguna Verde con vistas al mar.La cabaña cuenta con 3 dormitorios dos de ellos con camas de dos plazas y TV de pantalla plana. La Cocina se encuentra equipada con todo lo necesario para su utilización, contamos con nevera ,horno, campana y  todos los implementos para cocinar.Amplia zona de estar con una hermosa vista.Baño con agua caliente.Terraza , zona con parrilla y estacionamiento Privado.Viña del Mar se encuentra a tan solo 25 km de la casa.The spaceCabaña el Ocaso brinda un lugar para disfrutar de la naturaleza, vivir momentos de relax y disfrutar actividades al aire libre con una hermosa vista al mar.Guest accessLa cabaña se encuentra para uso exclusivo de los huéspedes, no se permiten visitas.Other things to noteLa cabaña tiene una capacidad máxima de 5 personas pero el cobro se realiza por huésped ( el precio  varía a mayor cantidad de huéspedes).No se permiten fiestas ni visitas.",
  "images": [
    "https://a0.muscache.com/im/pictures/271ff9e7-04c9-427e-b4ef-68b8135405b8.jpg",
    "https://a0.muscache.com/im/pictures/f10b8c0b-9ab5-4a92-8afa-5558d2a01145.jpg",
    "https://a0.muscache.com/im/pictures/80a0b6ce-9ed7-41c1-b9cd-10e28025f6c3.jpg",
    "https://a0.muscache.com/im/pictures/f8a49f5f-51a5-406f-a93f-83ea84032640.jpg",
    "https://a0.muscache.com/im/pictures/61446ecf-0fe8-46aa-a2f1-82e1b19d0865.jpg",
    "https://a0.muscache.com/im/pictures/71911216-1c87-46c4-b355-b563f2003055.jpg",
    "https://a0.muscache.com/im/pictures/fdb1b99d-a7d4-4fbb-8372-bd64ae73d3b5.jpg",
    "https://a0.muscache.com/im/pictures/c5777bcd-2fa6-4111-a37a-483a70de9d77.jpg"
  ]
}
```

En caso de que el room id sea inválido, el script devuelve un diccionario vacio:

```
# Ejecutando con un room id inválido
python main.py --room_id blahblablah
```

Output:

```
{}
```

Instrucciones para instalar airbnb-fetcher
==========================================

# Creación de un virtual env

1. Es requisito contar con Python 3 instalado. Preferentemente Python 3.8 en adelante.
2. Se recomienda crear un virtualenv para evitar conflictos con otros proyectos.

    Para crear un virtualenv:

    1. Instalar virtualenv

        ```shell
        pip3 install virtualenv
        ```

    2. Crear un virtualenv

        ```shell
        mkdir ~/.virtualenvs  # Podemos opcionalmente crear un directorio donde tener guardados los virtualenvs que creemos
        virtualenv -p python3.x ~/.virtualenvs/airbnb  # Reemplazar "python3.x" por la versión que tengamos
        ```

    3. Activar el virtualenv creado

        ```shell
        source ~/.virtualenv/airbnb/bin/activate
        ```

    4. Desactivar el virtualenv

        ```shell
        deactivate
        ```

3. Una vez dentros del virtualenv de `airbnb`, instalar las librerías necesarias con

    ```shell
    pip install -r requirements.txt
    ```

# Ejecutar airbnb-fetcher

1. Para ejecutar el script de airbnb manualmente, el primer paso es entrar el virtualenv de `airbnb`

    ```shell
    source ~/.virtualenv/airbnb/bin/activate
    ```

2. Ir al directorio del airbnb-fetcher

    ```shell
    cd airbnb-fetcher
    ```

3. Ejecutar pasando como room_id un room id válido de una propiedad de airbnb, como por ejemplo "34228829"

    ```shell
    python main.py --room_id 34228829
    ```

Instrucciones para ejecutar la imagen docker de prueba
======================================================

1. Instalar la imagen con

```
docker build . -t airbnb-fetcher:dev
```

2. Luego ejecutar la imagen con

```
docker run --rm airbnb-fetcher:dev --room_id <room_id>
```

Por ejemplo:

```
docker run --name airbnb-fetcher --rm airbnb-fetcher:dev --room_id 635437605163910390
```
