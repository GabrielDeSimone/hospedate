import utils from '../../utils/utils'

const backyardHost = utils.backyardHost
const backyardPort = utils.backyardPort

export default async function handler(req, res) {
  // Get data submitted in request's body.
  const body = req.body

  // Optional logging to see the responses
  // in the command line where next.js app is running.
  console.log('body: ', body)

  // Guard clause checks for first and last name,
  // and returns early if they are not found
  if (!body.title || !body.description || !body.price) {
    // Sends a HTTP bad request error code
    return res.status(400).json({ data: 'First or last name not found' })
  }

  const options = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({title: body.title, description: body.description, price: body.price}),
  }

  const response = await fetch(`http://${backyardHost}:${backyardPort}/posts`, options)

  const post = await response.json()

  // Found the name.
  // Sends a HTTP success code
  res.status(200).json({ data: `Cool with ${JSON.stringify(post)}` })
}