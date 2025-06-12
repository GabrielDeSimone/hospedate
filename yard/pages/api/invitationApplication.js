import { withIronSessionApiRoute } from "iron-session/next";
import utils from '../../utils/utils'

const apiErrors = utils.apiErrors;
const fetchBackyard = utils.fetchBackyard

async function invitationApplicationRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    const body = req.body

    if (!body.name || !body.email || !body.message) {
        return res.status(400).json({ error: 'name, email and message are required' })
    }
    if (body.name.length > utils.MAX_LENGTH_NAME_EXT_INVIT || body.email.length > utils.MAX_LENGTH_EMAIL ||
        body.message.length > utils.MAX_LENGTH_MESSAGE_EXT_INVIT) {
        return res.status(400).json({ error: 'Name email or message are too long' })
    }

    const [response, status] = await fetchBackyard('/notificationService/externalInvitationRequest', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            name: body.name,
            email: body.email,
            body: body.message,
        })
    }, res)

    if (status === 200) {
        res.status(200).json({
            data: "ok"
        })
    } else {
        res.status(500).json({
            data: apiErrors.UnknownError
        })
    }
}

export default withIronSessionApiRoute(invitationApplicationRoute, utils.ironOptions);
