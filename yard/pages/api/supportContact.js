import { withIronSessionApiRoute } from "iron-session/next";
import utils from '../../utils/utils'

const apiErrors = utils.apiErrors;

async function supportContactRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    res.status(200).json({
        data: {
            phoneNumber: utils.SUPPORT_PHONE_NUMBER,
            emailAddress: utils.SUPPORT_EMAIL_ADDRESS,
        }
    })
}

export default withIronSessionApiRoute(supportContactRoute, utils.ironOptions);
