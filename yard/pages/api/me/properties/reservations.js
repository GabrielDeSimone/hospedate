import utils from '../../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(PropsReservationsRoute, utils.ironOptions)

async function PropsReservationsRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const [response, status] = await fetchBackyard(`/orders/search?owner_id=${user.id}`)

    if (status === 200) {
        res.status(200).json({
            data: response.data.map((order) => ({
                id: order.id,
                propertyId: order.propertyId,
                status: order.status,
                checkinDate: order.date_start,
                checkoutDate: order.date_end,
                totalBilledCents: order.total_billed_cents,
                guests: order.number_guests,
                reservationType: order.order_type,
            }))
        })
    } else if (status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
