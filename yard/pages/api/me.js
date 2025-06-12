import { withIronSessionApiRoute } from "iron-session/next";
import utils from '../../utils/utils'

const apiErrors = utils.apiErrors;
const fetchBackyard = utils.fetchBackyard

async function meRoute(req, res) {
    if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }

    const user = req.session.user || null;

    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const [response, status] = await fetchBackyard(`/users/${user.id}`, null, res)

    const kindFromBackyard = {
        "for_traveler": utils.INVITATION_TYPE_FOR_GUEST,
        "for_owner": utils.INVITATION_TYPE_FOR_HOST,
    }

    if (status === 200) {
        const [responseInv, statusInv] = await fetchBackyard(`/invitations/search?generated_by=${user.id}`, null, res)
        if (statusInv === 200) {
            res.status(200).json({
                data: {
                    id: response.data.id,
                    name: response.data.name,
                    email: response.data.email,
                    createdAt: response.data.created_at,
                    phoneNumber: response.data.phone_number,
                    generatedInvitations: responseInv.data.map(invitation => ({
                        id: invitation.id,
                        kind: kindFromBackyard[invitation.kind],
                        used: Boolean(invitation.used_by),
                        createdAt: invitation.created_at
                    }))
                }
            })
        } else {
            res.status(500).json({ error: apiErrors.UnknownError })
        }
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}

export default withIronSessionApiRoute(meRoute, utils.ironOptions);
