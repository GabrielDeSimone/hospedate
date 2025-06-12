import utils from '../../utils/utils'
import { withIronSessionApiRoute } from "iron-session/next";

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

async function authenticateUser(req, user) {
    req.session.user = {
        id: user.id,
        name: user.name,
        email: user.email,
        isHost: user.is_host
    }
    await req.session.save();
}

export default withIronSessionApiRoute(loginRoute, utils.ironOptions)

async function loginRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const body = req.body

    if (!body.email || !body.password) {
        return res.status(400).json({ error: 'Email and password are required' })
    }
    if (body.email.length > utils.MAX_LENGTH_EMAIL || body.password.length > utils.MAX_LENGTH_PASSWORD) {
        return res.status(400).json({ error: 'Email or password are too long' })
    }

    const [response, status] = await fetchBackyard(`/users/search?email=${body.email}&password=${body.password}`, null, res)

    if (status === 200 && response.data.length > 0) {
        await authenticateUser(req, response.data[0])
        res.status(200).json({
            data: {
                user: {
                    id: response.data[0].id,
                    name: response.data[0].name,
                    email: response.data[0].email,
                    isHost: response.data[0].is_host,
                }
            }
        })
    } else if (status === 200) {
        res.status(401).json({ error: apiErrors.EmailOrPasswordIncorrect })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
