import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const backyardHost = utils.backyardHost
const backyardPort = utils.backyardPort
const apiErrors = utils.apiErrors

export default withIronSessionApiRoute(PropertiesSearchRoute, utils.ironOptions)

async function PropertiesSearchRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    // first we get the city, checkinDate, checkoutDate and guests from the query string
    const search = req.query
    const city = search.city
    const checkinDate = search.checkinDate
    const checkoutDate = search.checkoutDate
    const guests = search.guests

    // we create the options object for the fetch request
    const options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    }

    // we make the request to the backyard
    let propertiesSearch
    let response
    
    try {
        propertiesSearch = await fetch(
            [
                `http://${backyardHost}:${backyardPort}`,
                `/properties/search?city=${city}&date_start=${checkinDate}`,
                `&date_end=${checkoutDate}&status=active&guests=${guests}`
            ].join(''),
            options
        )
        response = await propertiesSearch.json()
    } catch (e) {
        console.log('Cannot connect to backyard', e)
        res.status(500).json({ error: apiErrors.UnknownError })
        return
    }

    if (propertiesSearch.status === 200) {
        res.status(200).json({
            data: response.data.map(property => ({
                id: property.id,
                title: property.title,
                maxGuests: property.max_guests,
                description: property.description,
                airbnb_room_id: property.airbnb_room_id,
                price: property.price,
                user_id: property.user_id,
                status: property.status,
                city: property.city,
                created_at: property.created_at,
                images: property.images,
            }))
        })
    } else if (propertiesSearch.status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
