import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard
const getOwnerId = utils.getOwnerId
const getOwner = utils.getOwner

export default withIronSessionApiRoute(ReservationsIdRoute, utils.ironOptions)

async function ReservationsIdRoute(req, res) {
    if (req.method !== "GET") {
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

    if (status === 200 && response.data) {

        const [responseUser, statusUser] = await fetchBackyard(`/users/${response.data.user_id}`, null, res)
        if (statusUser !== 200) {
            console.log(`Could not find user with id ${response.data.user_id} when enriching reservation`)
            res.status(500).json({ error: apiErrors.UnknownError })
            return
        }

        const ownerId = await getOwnerId(response.data.property_id, res)
        if (ownerId === user.id || response.data.user_id === user.id) {

            let owner = null;
            // if order is confirmed, add also owner info
            if (response.data.status === "confirmed") {
                owner = await getOwner(response.data.property_id)
            }

            res.status(200).json({
                data: {
                    id: response.data.id,
                    checkinDate: response.data.date_start,
                    checkoutDate: response.data.date_end,
                    totalBilledCents: response.data.total_billed_cents,
                    guests: response.data.number_guests,
                    propertyId: response.data.property_id,
                    ownerId: ownerId,
                    reservationType: response.data.order_type,
                    status: response.data.status,
                    walletAddress: response.data.wallet_address,
                    createdAt: response.data.created_at,
                    canceledBy: response.data.canceled_by,
                    user: {
                        id: response.data.user_id,
                        name: responseUser.data.name,
                        email: responseUser.data.email,
                        phoneNumber: responseUser.data.phone_number,
                    },
                    ...(owner && { owner: {
                            id: owner.id,
                            name: owner.name,
                            email: owner.email,
                    } })
                }
            })
        } else {
            res.status(404).json({ error: apiErrors.NotFound })
        }
    } else if (status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else if (status === 404) {
        res.status(404).json({ error: apiErrors.NotFound })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
