const backyardHost = process.env.BACKYARD_HOST
const backyardPort = process.env.BACKYARD_PORT

const SUPPORT_PHONE_NUMBER = "+1 786 940 6354"
const SUPPORT_EMAIL_ADDRESS = "info@hospedate.app"
const PROTECTED_RESERVATIONS_FEE_GUEST = 0.07
const PROTECTED_RESERVATIONS_FEE_HOST = 0.015
const MAX_LENGTH_EMAIL = 255
const MAX_LENGTH_PASSWORD = 255
const MAX_LENGTH_NAME_EXT_INVIT = 255
const MAX_LENGTH_MESSAGE_EXT_INVIT = 255

const INVITATION_TYPE_FOR_HOST = "FOR_HOST"
const INVITATION_TYPE_FOR_GUEST = "FOR_GUEST"

const INVITATION_MAX_LIMIT = 5

async function simpleFetchBackyard(relEndpoint, options) {
    if (!options) {
        options = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        }
    }

    try {
        const mainResponse = await fetch(`http://${backyardHost}:${backyardPort}${relEndpoint}`, options)
        const data = await mainResponse.json()
        return [data, mainResponse.status]
    } catch (e) {
        throw new Error(`Cannot connect to backyard ${e}`)
    }
}

async function getOwnerId(propertyId, res) {
    const [response, status] = await fetchBackyard(`/properties/${propertyId}`, null, res)
    if (status !== 200) {
        return null
    } else {
        return (response.data && response.data.user_id) || null
    }
}

async function getOwner(propertyId, res) {
    const ownerId = await getOwnerId(propertyId, res)
    if (ownerId) {
        const [response, status] = await fetchBackyard(`/users/${ownerId}`, null, res)
        if (status === 200 && response.data.id) {
            return response.data
        } else {
            return null
        }
    } else {
        return null
    }
}

async function fetchBackyard(relEndpoint, options, res) {
    try {
        const [data, status] = await simpleFetchBackyard(relEndpoint, options)
        return Promise.resolve([data, status])
    } catch (e) {
        res.status(500).json({ error: apiErrors.UnknownError })
        throw e;
    }
}

const apiErrors = {
    "EmailOrPhoneAlreadyExist": "EmailOrPhoneAlreadyExist",
    "EmailOrPasswordIncorrect": "EmailOrPasswordIncorrect",
    "PropertyAlreadyTaken": "PropertyAlreadyTaken",
    "PropertyActiveOrders": "PropertyActiveOrders",
    "InvitationNotValid": "InvitationNotValid",
    "BadRequest": "BadRequest",
    "NotFound": "NotFound",
    "UnknownError": "UnknownError",
    "UnsupportedMethod": "UnsupportedMethod",
    "Forbidden": "Forbidden"
}

const productVersion = process.env.NEXT_PUBLIC_PRODUCT_VERSION

const hospedateStringLogo = `

      ██    ╟█                                       j█▌         ]▄µ
      ██    ╟█  ,▓███▄  ,▓███▄  █▌▓██▓   ▄███▓   ▄██▓▄█▌ ,▓███▄ ▐███▓ _▄███▄
      ██▀▀▀▀██  █▌   ██ ╙█▄▄▄,  ██   ██ ▓█▄▄▄██ ╟█─  ╙█▌  ▄▄▄▓█▌ ▐█▌  ██▄▄▄██
      ██    ╟█  ██  ,██ ╓▄¬└╟█▌ ██  _██ ╟█▄  ▄▄ ╙█▄  ▄█▌ ██└ ▐█▌ ▐█▌  ██¬  ▄▄
      ▀▀    ╙▀   ▀▀▀▀▀   ▀▀▀▀▀  █▌▀▀▀▀   ╙▀▀▀▀─  ╙▀▀▀╙▀▀ ╙▀▀▀╙▀─  ▀▀▀  ╙▀▀▀▀
      ╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓ █▌ ╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓╓
      ${productVersion ? productVersion : ''}
`

const ironOptions = {
    cookieName: "hospedate_cookiename",
    password: "complex_password_at_least_32_characters_long",
    // secure: true should be used in production (HTTPS) but can't be used in development (HTTP)
    cookieOptions: {
        secure: process.env.NODE_ENV === "production"
    },
}

const cities = [
    "Buenos Aires",
    "Córdoba",
    "Rosario",
    "Mendoza",
    "Tucumán",
    "La Plata",
    "Mar del Plata",
    "Salta",
    "Santa Fe",
    "San Juan",
    "Resistencia",
    "Santiago del Estero",
    "Corrientes",
    "Neuquén",
    "Posadas",
    "Jujuy",
    "Bahía Blanca",
    "Paraná",
    "Formosa",
    "San Luis",
    "Catamarca",
    "Comodoro Rivadavia",
    "Río Cuarto",
    "Concordia",
    "San Nicolás de los Arroyos",
    "San Rafael",
    "Santa Rosa",
    "La Rioja",
]

const accommodations = [
    "house",
    "apartment",
    "private_room",
    "shared_room"
]

const locations = [
    "city_center",
    "near_beach",
    "residential_area",
    "countryside",
    "mountain",
]

const wifiOptions = [
    "shared",
    "private",
    "not_available",
]

const tvOptions = [
    "available",
    "available_cable_or_streaming",
    "not_available"
]

const parkingOptions = [
    "available_in_public_area",
    "available_private_uncovered",
    "available_private_covered",
    "not_available",
]

const propAmenitiesFieldsBackyardMap = {
    "accommodation": "accommodation",
    "location": "location",
    "wifi": "wifi",
    "tv": "tv",
    "microwave": "microwave",
    "oven": "oven",
    "kettle": "kettle",
    "toaster": "toaster",
    "coffeeMachine": "coffee_machine",
    "airConditioning": "air_conditioning",
    "heating": "heating",
    "parking": "parking",
    "pool": "pool",
    "gym": "gym",
    "halfBathrooms": "half_bathrooms",
    "bedrooms": "bedrooms",
}


export default {
    backyardHost,
    backyardPort,
    productVersion,
    apiErrors,
    ironOptions,
    hospedateStringLogo,
    fetchBackyard,
    getOwnerId,
    getOwner,
    cities,
    accommodations,
    locations,
    wifiOptions,
    tvOptions,
    parkingOptions,
    SUPPORT_PHONE_NUMBER,
    SUPPORT_EMAIL_ADDRESS,
    MAX_LENGTH_EMAIL,
    MAX_LENGTH_PASSWORD,
    MAX_LENGTH_NAME_EXT_INVIT,
    MAX_LENGTH_MESSAGE_EXT_INVIT,
    PROTECTED_RESERVATIONS_FEE_GUEST,
    PROTECTED_RESERVATIONS_FEE_HOST,
    propAmenitiesFieldsBackyardMap,
    INVITATION_TYPE_FOR_HOST,
    INVITATION_TYPE_FOR_GUEST,
    INVITATION_MAX_LIMIT,
}