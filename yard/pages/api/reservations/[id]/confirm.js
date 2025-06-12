import utils from '../../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard
const getOwnerId = utils.getOwnerId

export default withIronSessionApiRoute(ConfirmReservationRoute, utils.ironOptions)

async function ConfirmReservationRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    const user = req.session.user || null;
    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    if (!req.query.id) {
        res.status(400).json({ error: apiErrors.BadRequest })
        return
    }

    const [response, status] = await fetchBackyard(`/orders/${req.query.id}`, null, res)

    if (status === 200) {
        const ownerId = await getOwnerId(response.data.property_id, res)
        if (ownerId !== user.id) {
            res.status(403).json({ error: apiErrors.Forbidden })
        } else {
            const [responseEdit, statusEdit] = await fetchBackyard(`/orders/${req.query.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    "status": "confirmed",
                })
            }, res)
            if (statusEdit === 200) {
                res.status(200).json({
                    data: {
                        id: responseEdit.data.id,
                        checkinDate: responseEdit.data.date_start,
                        checkoutDate: responseEdit.data.date_end,
                        totalBilledCents: responseEdit.data.total_billed_cents,
                        guests: responseEdit.data.number_guests,
                        propertyId: responseEdit.data.property_id,
                        reservationType: responseEdit.data.order_type,
                        status: responseEdit.data.status,
                        walletAddress: responseEdit.data.wallet_address,
                        createdAt: responseEdit.data.created_at,
                    }
                })
            } else {
                console.log(`Problem confirming order: ${JSON.stringify((responseEdit))}`)
                res.status(500).json({ error: apiErrors.UnknownError })
            }
        }
    } else if (status === 404) {
        res.status(404).json({ error: apiErrors.NotFound })
    } else {
        console.log(`Unexpected status ${status} from backyard when fetching an order. Body was: ${JSON.stringify(response)}`)
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
