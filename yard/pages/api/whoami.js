import utils from '../../utils/utils'
import { withIronSessionApiRoute } from "iron-session/next";

const apiErrors = utils.apiErrors

export default withIronSessionApiRoute(whoAmIRoute, utils.ironOptions)

async function whoAmIRoute(req, res) {  
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null
    res.status(200).json({
        data: {
            user: user ? {
                id: user.id,
                name: user.name,
                email: user.email,
                isHost: user.isHost,
            } : null
        }
    })
}
