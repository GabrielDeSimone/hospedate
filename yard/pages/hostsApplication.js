import NavBar from '../components/navBar';
import {useContext} from "react";
import {UserContext, isGuest, isLoggedIn} from "../components/userProvider";
import Router from "next/router";
import ContentCard from "../components/contentCard";
import Link from "next/link";
import SmartContainer from "../components/smartContainer";

const MePage = () => {

    const [user, setUser] = useContext(UserContext);

    if (isGuest(user)) {
        Router.push('/')
    }


    return (
        <div>
            <NavBar />
            {isLoggedIn(user) && (
                <SmartContainer yAxisSpaced={true}>
                    <ContentCard>
                        <h1 className="text-3xl mb-4">Solitud para transformarte en anfitrión</h1>
                        <p>Por favor, contactate con el equipo de soporte para enviar tu solicitud. Podés encontrar el contacto en la sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help">Ayuda</Link>.</p>
                    </ContentCard>
                </SmartContainer>
            )}
        </div>
    )
}

export default MePage
