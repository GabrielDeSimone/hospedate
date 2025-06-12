import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const backyardHost = utils.backyardHost
const backyardPort = utils.backyardPort
const apiErrors = utils.apiErrors

export default withIronSessionApiRoute(MyPropertiesRoute, utils.ironOptions)

async function MyPropertiesRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    }

    let propertiesSearch
    let response
    
    try {
        propertiesSearch = await fetch(
            [
                `http://${backyardHost}:${backyardPort}`,
                `/properties/search?user_id=${user.id}`
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
            data: response.data.filter(property => property.status !== 'archived').map(property => ({
                id: property.id,
                title: property.title,
                description: property.description,
                maxGuests: property.max_guests,
                airbnb_room_id: property.airbnb_room_id,
                price: property.price,
                user_id: property.user_id,
                status: property.status,
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
