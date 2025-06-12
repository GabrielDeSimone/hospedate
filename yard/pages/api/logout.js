import { withIronSessionApiRoute } from "iron-session/next";
import utils from '../../utils/utils'

export default withIronSessionApiRoute(
  function logoutRoute(req, res, session) {

    req.session.destroy();
    res.json({ message: 'success' });

  }, utils.ironOptions);
