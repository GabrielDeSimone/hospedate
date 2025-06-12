import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'
const fetchBackyard = utils.fetchBackyard

const apiErrors = utils.apiErrors

export default withIronSessionApiRoute(MyReservationsRoute, utils.ironOptions)

async function MyReservationsRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const [response, status] = await fetchBackyard(`/orders/search?user_id=${user.id}`, null, res)

    if ([200, 404].includes(status)) {
        res.status(status)
    } else {
        res.status(500)
        console.log("Error while calling register, got this from backend: " + JSON.stringify(response))
    }

    if (status === 200) {
        res.json({ data: response.data.map(order => ({
            id: order.id,
            propertyId: order.propertyId,
            status: order.status,
            checkinDate: order.date_start,
            checkoutDate: order.date_end,
            totalBilledCents: order.total_billed_cents,
            guests: order.number_guests,
            reservationType: order.order_type,
        }))})
    } else if (status === 400) {
        res.json({ error: apiErrors.BadRequest })
    } else {
        res.json({ error: apiErrors.UnknownError })
    }
}
