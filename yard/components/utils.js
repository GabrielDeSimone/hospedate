function humanizeDate(dateStr, options) {
    if (!options) {
        options = { day: 'numeric', month: 'long', timeZone: 'UTC' };
    }
    const [year, month, day] = dateStr.split('-').map(Number);
    const date = new Date(Date.UTC(year, month - 1, day));
    return date.toLocaleDateString('es-ES', options);
}

function humanizeCents(amountCents) {
    return Number(amountCents / 100).toFixed(2)
}

function humanizeTimestamp(timestamp) {
    const date = new Date(timestamp);

    // Extracting individual components
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0'); // Months are 0-based
    const day = date.getDate().toString().padStart(2, '0');
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

function humanizeBoolean(bool) {
    return bool ? "Sí" : "No"
}

function humanizeInvitationKind(invKind) {
    return {
        "FOR_HOST": "Anfitrión",
        "FOR_GUEST": "Huésped",
    }[invKind]
}

function humanizeFullDate(dateStr) {
    return humanizeDate(dateStr, { day: 'numeric', month: 'long', timeZone: 'UTC', year: 'numeric' })
}

function getTodayDate() {
    return (new Date()).toISOString().slice(0,10)
}

function getLocalTodayDate() {
    return humanizeTimestamp(new Date()).slice(0,10)
}

function getDaysDifference(checkinDate, checkoutDate) {
    let date1 = new Date(checkinDate);
    let date2 = new Date(checkoutDate);

    let differenceInTime = date2.getTime() - date1.getTime();

    return Math.floor(differenceInTime / (1000 * 3600 * 24));
}

const RESERVATION_STATUS_EPHEMERAL = "ephemeral"
const RESERVATION_STATUS_DISCARDED = "discarded"
const RESERVATION_STATUS_PENDING = "pending"
const RESERVATION_STATUS_CONFIRMED = "confirmed"
const RESERVATION_STATUS_CANCELED = "canceled"
const RESERVATION_STATUS_IN_PROGRESS = "in_progress"
const RESERVATION_STATUS_COMPLETED = "completed"

const CANCELED_BY_OWNER = "owner"
const CANCELED_BY_VISITOR = "visitor"

function humanizeReservationStatus(status) {
    if (status === RESERVATION_STATUS_EPHEMERAL) {
        return "Pendiente de pago"
    } else if (status === RESERVATION_STATUS_DISCARDED) {
        return "Descartada"
    } else if (status === RESERVATION_STATUS_PENDING) {
        return "Pendiente"
    } else if (status === RESERVATION_STATUS_CONFIRMED) {
        return "Confirmada"
    } else if (status === RESERVATION_STATUS_CANCELED) {
        return "Cancelada"
    } else if (status === RESERVATION_STATUS_COMPLETED) {
        return "Completada"
    } else {
        return "-"
    }
}

const humanizeReservationType = reservationType => {
    if (reservationType === RESERVATION_HOSPEDATE) {
        return "Reserva protegida"
    } else if (reservationType === RESERVATION_DIRECT_OWNER) {
        return "Reserva directa"
    } else {
        return "-"
    }
}

function remainingEphemeralSeconds(createdAt) {
    return Math.round(3600 - new Date().getTime()/1000 + new Date(createdAt).getTime()/1000)
}

function humanizeRemainingTime(remainingSecs) {
    function humanizeMins(remSecs) {
        const mins = Math.round(remSecs/60)
        return `${mins} minuto${mins === 1? "" : "s"}`
    }
    const hours = Math.floor(remainingSecs/3600)
    if (hours > 0) {
        return `${hours} hora${hours === 1? "" : "s"} ${humanizeMins(remainingSecs-hours*3600)}`
    } else {
        return humanizeMins(remainingSecs)
    }
}

function reservationStatusIsFinal(status) {
    return status === RESERVATION_STATUS_CONFIRMED || status === RESERVATION_STATUS_DISCARDED
}

const amenitiesList = [
    {
        icon: 'wc',
        key: 'halfBathrooms',
        strValue: (x) => `${x / 2}`,
        strFull: (x) => `${x / 2} baños`,
        strText: "Baños"
    },
    {
        icon: 'bed',
        key: 'bedrooms',
        strValue: (x) => `${x}`,
        strFull: (x) => `${x} dormitorios`,
        strText: "Dormitorios"
    },
    {
        icon: 'home',
        key: 'accommodation',
        strValue: (x) => `${humanizeAccommodation(x)}`,
        strFull: (x) => `${humanizeAccommodation(x)}`,
        strText: "Tipo de alojamiento"
    },
    {
        icon: 'map',
        key: 'location',
        strValue: (x) => `${humanizeLocation(x)}`,
        strFull: (x) => `${humanizeLocation(x)}`,
        strText: "Ubicación"
    },
    {
        icon: 'wifi',
        key: 'wifi',
        strValue: (x) => `${humanizeWifiOption(x)}`,
        strFull: (x) => `${humanizeWifiOption(x)}`,
        strText: "Opción de Wifi"
    },
    {
        icon: 'tv',
        key: 'tv',
        strValue: (x) => `${humanizeTvOption(x)}`,
        strFull: (x) => `${humanizeTvOption(x)}`,
        strText: "Opción de Tv"
    },
    {
        icon: 'garage_home',
        key: 'parking',
        strValue: (x) => `${humanizeParkingOption(x)}`,
        strFull: (x) => `${humanizeParkingOption(x)}`,
        strText: "Estacionamiento"
    },
    {
        icon: 'microwave',
        key: 'microwave',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Microondas",
        isBoolean: true,
    },
    {
        icon: 'oven_gen',
        key: 'oven',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Horno",
        isBoolean: true,
    },
    {
        icon: 'kettle',
        key: 'kettle',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Pava eléctrica",
        isBoolean: true,
    },
    {
        icon: 'breakfast_dining',
        key: 'toaster',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Tostadora",
        isBoolean: true,
    },
    {
        icon: 'coffee_maker',
        key: 'coffeeMachine',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Máquina de café",
        isBoolean: true,
    },
    {
        icon: 'mode_fan',
        key: 'airConditioning',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Aire acondicionado",
        isBoolean: true,
    },
    {
        icon: 'fireplace',
        key: 'heating',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Calefacción",
        isBoolean: true,
    },
    {
        icon: 'pool',
        key: 'pool',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Piscina",
        isBoolean: true,
    },
    {
        icon: 'exercise',
        key: 'gym',
        strValue: (x) => `${humanizeBoolean(x)}`,
        strFull: (x) => `${humanizeBoolean(x)}`,
        strText: "Gimnasio",
        isBoolean: true,
    },
]

function humanizeAccommodation(acc) {
    return ({
        "house": "Casa",
        "apartment": "Departamento",
        "private_room": "Habitación privada",
        "shared_room": "Habitación compartida"
    }[acc])
}

function humanizeLocation(loc) {
    return ({
        "city_center": "Centro de la ciudad",
        "near_beach": "Cerca de la playa",
        "residential_area": "Área residencial",
        "countryside": "Área rural",
        "mountain": "Montaña",
    }[loc])
}

function humanizeWifiOption(wifi) {
    return ({
        "shared": "Compartido",
        "private": "Privado",
        "not_available": "No disponible"
    }[wifi])
}

function humanizeTvOption(tvOpt) {
    return ({
        "available": "Disponible",
        "available_cable_or_streaming": "Tv con servicio de cable o streaming",
        "not_available": "No disponible",
    }[tvOpt])
}

function humanizeParkingOption(parkingOpt) {
    return ({
        "available_in_public_area": "Estacionamiento en lugar público",
        "available_private_uncovered": "Estacionamiento en lugar privado",
        "available_private_covered": "Estacionamiento privado y techado",
        "not_available": "No disponible",
    }[parkingOpt])
}


const RESERVATION_HOSPEDATE = "in_platform"
const RESERVATION_DIRECT_OWNER = "owner_directly"

const utils = {
    humanizeDate,
    getDaysDifference,
    humanizeFullDate,
    humanizeTimestamp,
    getTodayDate,
    getLocalTodayDate,
    humanizeBoolean,
    humanizeCents,
    RESERVATION_HOSPEDATE,
    RESERVATION_DIRECT_OWNER,
    humanizeReservationStatus,
    humanizeReservationType,
    RESERVATION_STATUS_EPHEMERAL,
    RESERVATION_STATUS_PENDING,
    RESERVATION_STATUS_CONFIRMED,
    RESERVATION_STATUS_IN_PROGRESS,
    RESERVATION_STATUS_COMPLETED,
    RESERVATION_STATUS_CANCELED,
    RESERVATION_STATUS_DISCARDED,
    CANCELED_BY_OWNER,
    CANCELED_BY_VISITOR,
    remainingEphemeralSeconds,
    humanizeRemainingTime,
    reservationStatusIsFinal,
    humanizeAccommodation,
    humanizeLocation,
    humanizeWifiOption,
    humanizeTvOption,
    humanizeParkingOption,
    humanizeInvitationKind,
    amenitiesList,
}


export default utils