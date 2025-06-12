import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(MyCreditRoute, utils.ironOptions)

async function MyCreditRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const [response, status] = await fetchBackyard(`/users/${user.id}/credit`, null, res)

    if (status === 200) {
        res.status(200).json({
            data: response.data.map((userCredit) => ({
                id: userCredit.id,
                invitationId: userCredit.invitation_id,
                earnedAmount: userCredit.earned_amount,
                earnedCurrency: userCredit.earned_currency,
                createdAt: userCredit.created_at,
            }))
        })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
