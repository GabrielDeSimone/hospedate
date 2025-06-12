import { withIronSessionApiRoute } from "iron-session/next";
import utils from '../../utils/utils'

const apiErrors = utils.apiErrors;
const backyardHost = utils.backyardHost
const backyardPort = utils.backyardPort

async function postsRoute(req, res) {
    const postsResp = await fetch(`http://${backyardHost}:${backyardPort}/posts`)
    const response = await postsResp.json()

    if (response.data) {
        res.status(200).json({ data: response.data })
    } else {
        if (response.error) {
            console.log(`Error when fetching posts: ${response.error}`);
        }
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}

export default withIronSessionApiRoute(postsRoute, utils.ironOptions);
