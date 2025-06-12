import {useEffect, useState} from "react";
import HospedateLogo from "./hospedateLogo";
import FormButton from "./formButton";
import FormInputIcon from "./formInputIcon";
import toast from "react-hot-toast";
import Router from "next/router";


const INVITATION_SCREEN_ID = "landingScreen6"

const WelcomeHosts = () => {

    const [whiteScreenHeight, setWhiteScreenHeight] = useState("h-screen")
    const [logoFlexDirection, setLogoFlexDirection] = useState("flex-col")
    const [dummyExpanderClasses, setDummyExpanderClasses] = useState("flex-initial")
    const [logoHeight, setLogoHeight] = useState(50)
    const [inProgressChatResponse, setInProgressChatResponse] = useState("")
    const [chatResponse, setChatResponse] = useState("")
    const [showChatResponse, setShowChatResponse] = useState(false)
    const [showIdentityForm, setShowIdentityForm] = useState(false)
    const [chatFullExpanded, setChatFullExpanded] = useState(false)

    async function handleChatSend(e) {
        e.preventDefault()
        if (showChatResponse || inProgressChatResponse.trim() === "") {
            return
        }
        setChatResponse(inProgressChatResponse)
        setInProgressChatResponse("")
        setShowChatResponse(true)
        setTimeout(() => {
            setShowIdentityForm(true)
        }, 300)
    }

    async function handleWantInvitation(e) {
        e.preventDefault()
        const element = document.getElementById(INVITATION_SCREEN_ID)
        if (element) {
            element.scrollIntoView({ behavior: "smooth" });
        }
    }

    async function handleHaveInvitation(e) {
        e.preventDefault()
        Router.push("/register")
    }

    async function handleIdentitySubmit(event) {
        event.preventDefault()
        const name = event.target.name.value.trim()
        const email = event.target.email.value.trim()

        if (!name || !email) {
            return
        }

        const response = await fetch("/api/invitationApplication", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name,
                email,
                message: chatResponse,
            })
        })
        const result = await response.json()
        if (result.error) {
            toast.error("Hubo un error al mandar el mensaje")
        } else {
            setShowIdentityForm(false)
            setTimeout(() => {
                setChatFullExpanded(true)
            }, 500)
        }
    }

    useEffect(() => {
        setTimeout(() => {
            setWhiteScreenHeight('h-[84px]')
            setLogoHeight(null)
            setTimeout(() => {
                setLogoFlexDirection("flex-row items-center")
                setDummyExpanderClasses("flex-auto")
            }, 300)
        }, 400)
    }, [])

    return (
        <div>
            <div className={`myAwesomeClass px-[50px] min-[500px]:px-[100px] min-[1240px]:px-[200px] mx-auto flex justify-center ${logoFlexDirection} ${whiteScreenHeight}`}>
                <HospedateLogo height={logoHeight} />
                <div className={dummyExpanderClasses + " myAwesomeClass2"} />
            </div>


            <div className="landingScreen landingScreen1">
                {/*text-[90px] mb-5 text-white font-bold text-center px-5*/}
                <h1 className="">Descubrí un <br className="tablet:hidden" /> nuevo <br className="hidden tablet:block" />ingreso <br className="tablet:hidden" /> para tu propiedad</h1>
                <p>Más de <span>9.3 millones</span> de usuarios<br />de criptomonedas en Argentina<br />y <span>420 millones</span> en todo el mundo<br />te están esperando.</p>
                <div>
                    <form className="mt-7 px-5 tablet:px-0 flex flex-col tablet:flex-row justify-center">
                        <FormButton onClick={handleWantInvitation} smRounded extraClasses="h-[60px] pt-[15px] border transition" bgClasses="bg-hospedate-green hover:bg-opacity-70" icon="key" text="Quiero una invitación" />
                        <FormButton onClick={handleHaveInvitation} smRounded extraClasses="h-[60px] pt-[15px] tablet:ml-3 border transition" bgClasses="bg-black hover:bg-opacity-70" icon="hand_gesture" text="Tengo una invitación" />
                    </form>
                </div>
            </div>

            <div className="landingScreen landingScreen2 laptop:px-32">
                <h2>No te ajustes a<br />una plataforma,<br /><span>nosotros nos<br />ajustamos a vos</span></h2>
                <p>Arreglá el pago de forma directa con<br />tus huéspedes o usá nuestro sistema<br />de pagos, <span>la elección es tuya.</span></p>
            </div>

            <div className="landingScreen landingScreen3 textRight laptop:px-32">
                <h2>Inconfiscable,<br />ultra rápido y<br /><span>sin burocracia</span></h2>
                <p>Recibí tu dinero en <span>USDT</span><br className="tablet:hidden" /> en tu <br className="hidden tablet:block" />billetera de confianza. <br className="tablet:hidden" />Nosotros <br className="hidden tablet:block" /> te asesoramos <br className="tablet:hidden" /> en cada paso.</p>
            </div>

            <div className="landingScreen landingScreen4 laptop:px-32">
                <h2>Llená tu calendario,<br /><span>fácil y sin vueltas</span></h2>
                <p>Ahorrá tiempo usando <br className="tablet:hidden" /><span>OneClick Mirror</span> para <br />importar tu propiedad <br className="tablet:hidden" /> desde otras plataformas.<br /><span>Tan simple como copiar <br className="tablet:hidden" /> y pegar un link.</span></p>
            </div>

            <div className="landingScreen landingScreen5 textRight textTop laptop:px-32">
                <h2>Tus mejores <br className="tablet:hidden" /> aliados en<br /><span>un mercado <br className="tablet:hidden" /> desafiante</span></h2>
                <p>Entendemos la complejidad de <br className="tablet:hidden" /> ser propietario hoy en día.<br />Nuestra dedicación y atención <br className="tablet:hidden" /> al cliente es la garantía de que<br /><span>tu éxito es nuestro éxito.</span></p>
            </div>

            <div id={INVITATION_SCREEN_ID} className="landingScreen landingScreen6 laptop:px-32">
                <h2>¡Nos interesa <span>tu opinión!</span></h2>
                <div className="chat w-full tablet:w-[600px]">
                    <div className={`identityPopup ${showIdentityForm ? "block" : "hidden"}`}>
                        <form onSubmit={handleIdentitySubmit} >
                            <p>Completá tus datos para que podamos recibir tu mensaje.</p>
                            <FormInputIcon icon="badge" name="name" id="name" placeholder="Nombre" required />
                            <FormInputIcon icon="email" type="email" name="email" id="email" placeholder="Correo electrónico" required />
                            <FormButton type="submit" text="Enviar" />
                        </form>
                    </div>
                    <div className={`messages ${chatFullExpanded ? "fullExpanded" : (showChatResponse ? "expanded" : "")}`}>
                        <div>
                            <p>Hola! Estamos en proceso de beta. Si te interesa formar parte de los pioneros y recibir beneficios exclusivos, envianos tu consulta por acá.</p>
                            <img src="solapaizq.svg" />
                        </div>
                        <div className="response">
                            <p>{chatResponse}</p>
                            <img src="solapader.svg" />
                        </div>
                        <div className="finalResponse">
                            <p>Gracias por escribirnos. Pronto nos pondremos en contacto con vos!</p>
                            <img src="solapaizq.svg" />
                        </div>
                    </div>
                    <form onSubmit={handleChatSend} className="controls">
                        <input type="text" disabled={showChatResponse} placeholder="Mensaje" value={inProgressChatResponse} onChange={(e) => setInProgressChatResponse(e.target.value)} />
                        <button onClick={handleChatSend}>
                            {/*<span className="material-symbols-outlined">send</span>*/}
                            <img src="/sendIcon.svg" />
                        </button>
                    </form>
                </div>
            </div>
            {/*<Footer />*/}
        </div>
    )
}

export default WelcomeHosts
