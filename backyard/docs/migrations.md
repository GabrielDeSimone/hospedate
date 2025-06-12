
Definición de una migración y del conjunto de migraciones
==============

Una migración es una especie de "update" que le aplicamos a la BD para llevarla de una versión más antigua a una más moderna. El ejemplo típico es agregar un nuevo atributo a una tabla. El requerimiento es siempre preservar los datos existentes en una tabla, excepto que una migración particular no lo requiera.

Definimos una migración como una terna con 3 elementos

(Va, Vb, [q1,...,qn])

Donde

Va: es la "version A". Representa el estado inicial de la BD

Vb: es la "versión B". Representa el estado final de la BD

[q1,...,qn]: Es un listado de queries que de ser ejecutadas exitosamente en la BD_Va, la transforman en BD_Vb

Además, haciendo una analogía con espacios vectoriales - con un espacio vectorial para cada versión de la BD y los vectores del mismo dados por conjuntos de datos posibles en esa versión -, las migraciones son transformaciones que consideramos inversibles. Esto es, que existe un conjunto de queries que permiten volver de BD_Vb a BD_Va. Esta operación de inversión (que no necesariamente es única), como no podemos calcularla, la agregamos como un elemento adicional de la definición de migración. De este modo:

(Va, Vb, [q1,...,qn],[qi_1,...,qi_m])

Considerando un conjunto de versiones posibles (o espacios vectoriales) y un conjunto de migraciones como formas de "moverse" de una versión a otra (transformaciones), es posible representar ambas como un grafo. En este sentido, es posible tener muchas tipologías. Si se impone los siguientes axiomas de validez a los conjuntos de migraciones:
- Cada versión solo puede ir a una única versión
- Existe una única versión de origen (no partimos de múltiples versiones posibles)
Se sigue que, para ser válido, un conjunto de migraciones tiene que constituir una cadena completa (sin cortes) y en la que cada versión no puede aparecer más de una vez como versión de origen o como versión de destino en el conjunto. 
El motivo en la práctica para imponer los axiomas anteriores es que facilita la aplicación secuencial de migraciones a una base de datos. Estos conceptos se aplican para gestionar y mantener bases de datos en evolución, asegurando que los datos existentes se conserven durante las actualizaciones y permitiendo la reversión de cambios si es necesario.


Migraciones en Go
==============

Las migraciones y el conjunto de migraciones en go los definimos a partir de los siguientes tipos de datos que se encuentran en globalRepository

type dbversion string

type migration struct {
    VersionFrom   dbversion
    VersionTo     dbversion
    Queries       []string
    QueriesInvert []string
}

type MigrationMap map[dbversion]*migration

Se presenta seguidamente un ejemplo de cómo quedaría implementado un conjunto de migraciones con una sola migración. En el caso de que no hayan migraciones, se deja una lista vacía. 
La lista se define en el fichero migrationsList.go y las queries asociadas a cada migración en migrations.go

```
var migrations = []migration{
    {
        VersionFrom:   INITIAL_VERSION,
        VersionTo:     dbversion("v1"),
        Queries:       sqlTmp.QueriesV0ToV1,
        QueriesInvert: sqlTmp.QueriesInvertV0ToV1,
    },
}

var QueriesV0ToV1 = []string{
    "ALTER TABLE users ADD age int;",
}

var QueriesInvertV0ToV1 = []string{
    "ALTER TABLE users DROP COLUMN age;",
}
```
