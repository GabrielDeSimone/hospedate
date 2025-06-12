import utils from '../../utils/utils'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default async function handler(req, res) {
    const body = req.body

    if (!body.email || !body.password || !body.name || !body.phoneNumber || !body.invitationId) {
        res.status(400).json({ error: 'Email, password, name, phoneNumber and invitationId are required' })
        return
    }

    const options = {
        method: 'POST',
        headers: {
        'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: body.email,
            password: body.password,
            name: body.name,
            phone_number: body.phoneNumber,
            invitation_id: body.invitationId ? body.invitationId : "", // we DON'T want to send a null accidentally
        }),
    }

    const [response, status] = await fetchBackyard('/users', options, res)

    if ([200, 201, 400].includes(status)) {
        res.status(status)
    } else {
        res.status(500)
        console.log("Error while calling register, got this from backend: " + JSON.stringify(response))
    }

    if (String(status)[0] === '2') {
        res.json({ data: { id: response.id } })
    } else if (status === 400 && response.error === "ErrEmailOrPhoneAlreadyTaken") {
        res.json({ error: apiErrors.EmailOrPhoneAlreadyExist })
    } else if (status === 400 && response.error === "ErrInvitationNotValid") {
        res.json({ error: apiErrors.InvitationNotValid })
    } else if (status === 400 && response.error === "ErrInvitationDoesNotExist") {
        res.json({ error: apiErrors.InvitationNotValid })
    } else if (status === 400 && response.error === "ErrInvitationAlreadyUsed") {
        res.json({ error: apiErrors.InvitationNotValid })
    } else if (status === 400) {
        res.json({ error: apiErrors.BadRequest })
    } else {
        res.json({ error: apiErrors.UnknownError })
    }
}