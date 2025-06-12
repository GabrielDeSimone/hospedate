Instrucciones para ejecutar los integration tests
=================================================


# Instalación de behave

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
        virtualenv -p python3.x ~/.virtualenvs/backyard-tests  # Reemplazar "python3.x" por la versión que tengamos
        ```

    3. Activar el virtualenv creado

        ```shell
        source ~/.virtualenv/backyard-tests/bin/activate
        ```

    4. Desactivar el virtualenv

        ```shell
        deactivate
        ```

3. Una vez dentros del virtualenv de `backyard-tests`, instalar las librerías necesarias con

    ```shell
    pip install -r requirements.txt
    ```

# Ejecutar los tests

1. Para ejecutar los tests, el primer paso es entrar el virtualenv de `backyard-tests`

    ```shell
    source ~/.virtualenv/backyard-tests/bin/activate
    ```

2. Ir al directorio de integration tests

    ```shell
    cd integration_tests
    ```

3. Ejecutar behave

    ```shell
    behave
    ```

# Ejecutar un test en particular

Para ejecutar un archivo de features en particular:

```shell
behave -i features/users.feature
```

Para ejecutar un escenario en particular:

```shell
behave -n "Search a property by city and dates"
```

Tener en cuenta que esto va a ejecutar todos los escenarios que contengan el string "Search a property by city and dates" por lo que si el nombre del escenario es bastante simple como "Create an user" es posible que más de un test sea ejecutado.

Para ejecutar test que contengan ciertos TAGS:

```shell
behave -t @MYTAG
```

Esta expresión va a ejecutar escenarios que contengan el TAG @MYTAG

Se recomienda también usar el flag `--no-skipped` o `-k` para evitar ver logs de tests no ejecutados

```shell
behave -t @MYTAG -k
```