import utils from '../../utils/utils'
import { withIronSessionApiRoute } from "iron-session/next";

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(invitationsRoute, utils.ironOptions)

const godUserIds = ["rVXGSAIc4P", "cvjkbsURCK", "ba03VufCig"]


async function currentNumberInvitations(userId, res) {
    const [responseInv, statusInv] = await fetchBackyard(`/invitations/search?generated_by=${userId}`, null, res)
    if (statusInv === 200) {
        return responseInv.data.length
    } else {
        return null
    }
}


async function invitationsRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    if (!req.session.user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    const user = req.session.user
    const body = req.body
    if (!body.kind) {
        res.status(400).json({ error: 'kind is required' })
        return
    }
    if (body.kind !== utils.INVITATION_TYPE_FOR_GUEST && body.kind !== utils.INVITATION_TYPE_FOR_HOST) {
        res.status(400).json({ error: 'invalid value for kind' })
        return
    }

    const numInvitationsGenerated = await currentNumberInvitations(user.id, res)
    if (!godUserIds.includes(user.id) && numInvitationsGenerated >= utils.INVITATION_MAX_LIMIT) {
        res.status(400).json({ error: 'Maximum invitations limit' })
        return
    }

    const kindToBackyard = {
        [utils.INVITATION_TYPE_FOR_GUEST]: "for_traveler",
        [utils.INVITATION_TYPE_FOR_HOST]: "for_owner",
    }

    const kindFromBackyard = {
        "for_traveler": utils.INVITATION_TYPE_FOR_GUEST,
        "for_owner": utils.INVITATION_TYPE_FOR_HOST,
    }

    const [response, status] = await fetchBackyard('/invitations', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            kind: kindToBackyard[body.kind],
            generated_by: user.id
        })
    }, res)

    if (status === 201) {
        res.status(201).json({
            data: {
                id: response.data.id,
                kind: kindFromBackyard[response.data.kind],
                used: Boolean(response.data.used_by),
                createdAt: response.data.created_at
            }
        })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
