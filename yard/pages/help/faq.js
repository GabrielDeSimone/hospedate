import NavBar from '../../components/navBar';
import {useContext, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../../components/userProvider";
import Router from "next/router";
import ContentCard from "../../components/contentCard";
import Link from "next/link";
import SmartContainer from "../../components/smartContainer";


const initialFaqInfo = {
    sections: [
        {
            title: "Registro e invitaciones",
            questions: [
                {
                    question: "¿Cómo puedo registrarme en Hospedate?",
                    answer: "Por el momento, para registrarse en Hospedate es necesario tener una invitación. Como usuario podés generar invitaciones desde la sección \"Mi cuenta\" para que otras personas se registren.",
                },
                {
                    question: "¿Qué tipos de invitaciones existen?",
                    answer: "Existen dos tipos de invitaciones. Invitaciones para huéspedes e invitaciones para anfitriones. Los usuarios registrados mediante una invitación para huéspedes también pueden aplicar para que su cuenta como anfitrión sea habilitada.",
                },
            ]
        },
        {
            title: "Acerca de Hospedate",
            questions: [
                {
                    question: "¿Qué es Hospedate?",
                    answer: "Hospedate es una plataforma de alquileres temporarios que contacta anfitriones de propiedades con huéspedes. Hospedate ofrece un servicio de Reservas Protegidas a través de criptomonedas y también Reservas Directas, donde las partes se contactan directamente y deciden la forma de pago que les es más conveniente según sus preferencias.",
                },
                {
                    question: "¿Qué tiene Hospedate de diferente con otras plataformas de hospedajes que ya existen?",
                    answer: "A diferencia de otras plataformas, Hospedate permite a sus usuarios optar entre Reservas Protegidas o Reservas Directas. Los usuarios son libres de contactarse y elegir el método de pago más conveniente. Además, a diferencia de otras plataformas, en Hospedate podés hacer reservas con criptomonedas.",
                }
            ]
        },
        {
            title: "Búsqueda y carga de propiedades",
            questions: [
                {
                    question: "¿Cómo puedo ser anfitrión?",
                    answer: "Para aplicar como anfitrión, debés contactar al equipo de soporte desde la sección de Ayuda. El equipo dará de alta tu cuenta como anfitrión y podrás cargar propiedades.",
                },
                {
                    question: "¿Cómo puedo cargar mi propiedad en Hospedate?",
                    answer: "Para cargar tu propiedad, debés ir a la sección \"Anfitriones\". Desde allí hacer click en \"Agregar nueva propiedad\".",
                },
                {
                    question: "¿Cómo se realiza una reserva?",
                    answer: "Por el momento, un anfitrión debe compartir el link de su propiedad a un huésped interesado en reservar. Si el huésped no tiene cuenta en Hospedate, entonces deberá registrarse con una invitación válida.",
                },
                {
                    question: "¿Cómo puedo buscar propiedades?",
                    answer: "Por el momento, la búsqueda de propiedades está deshabilitada hasta que la plataforma cuente con una cantidad significativa de propiedades cargadas. ",
                },
            ]
        },
        {
            title: "Reservas",
            questions: [
                {
                    question: "¿Qué tipos de pagos se pueden hacer desde Hospedate?",
                    answer: "Hospedate acepta reservas de dos tipos: Reservas Protegidas o Reservas Directas. Las reservas protegidas se pagan con criptomonedas. En las reservas directas ambas partes coordinan el pago de la forma más conveniente."
                },
                {
                    question: "¿Puedo pagar con tarjeta de crédito o débito?",
                    answer: "Si la reserva es Protegida, entonces no es posible pagar con tarjetas tradicionales. Sin embargo, si la reserva es de tipo Directa, entonces podrás pagar con cualquier método de pago dispuesto por el dueño."
                },
                {
                    question: "Como anfitrión, ¿cómo me pagan cuando un huésped hace una Reserva Protegida?",
                    answer: "Una vez que el check-in se ha completado exitosamente, Hospedate realiza una acreditación de fondos que podés ver desde el panel de Anfitriones. Una vez que tu balance esté disponible, podés hacer un retiro a una wallet que sea de tu custodia."
                },
                {
                    question: "¿Cómo funciona la política de cancelaciones de Hospedate?",
                    answer: "Podés consultar nuestra Política de Cancelaciones en la sección de Ayuda."
                },
            ]
        },
        {
            title: "Pagos y criptomonedas",
            questions: [
                {
                    question: "¿Cómo hago un retiro de fondos a una wallet propia?",
                    answer: "Actualmente debés contactarnos por WhatsApp al equipo de soporte, cuyos datos de contacto podrás encontrar en la sección de Ayuda. Estamos trabajando en una opción de retiros automatizada para tu comodidad."
                },
                {
                    question: "¿Qué es una wallet?",
                    answer: "Una wallet es una aplicación donde podés tener custodia de tus criptomonedas. Desde ahí, podés enviar tu dinero a terceros o realizar compras de bienes y servicios con tus criptomonedas."
                },
                {
                    question: "¿Qué tipos de wallets existen?",
                    answer: "Principalmente existen las wallets custodiales y las no custodiales. En las wallets custodiales, un tercero de confianza (un custodio) tiene el control y la responsabilidad de almacenar y proteger las claves privadas asociadas con las criptomonedas de un usuario. En las wallets no custodiales, el usuario tiene el control total sobre las claves privadas. Esto significa que el usuario es el único responsable de almacenar y proteger sus claves privadas. "
                },
                {
                    question: "¿Qué wallets me conviene usar?",
                    answer: "Existen miles de wallets donde podés guardar tus USDT. Entre las wallets custodiales, existen los exchanges, que actúan como bancos, como Binance o Kucoin. Un ejemplo de wallet no custodial es TrustWallet. También, existen las wallets de hardware, que son dispositivos físicos diseñados específicamente para el almacenamiento seguro de criptomonedas, como Ledger o Trazor. Son altamente seguros y a menudo se consideran la opción más segura. En cualquier caso, te recomendamos investigar los pros y contras de las wallets disponibles en el mercado para que tomes la decisión que te deje más tranquilo."
                },
                {
                    question: "¿Qué criptomonedas utiliza Hospedate?",
                    answer: "Actualmente, Hospedate utiliza USDT como criptomoneda para Reservas Protegidas."
                },
            ]
        },
        {
            title: "Comisiones y cargos de Hospedate",
            questions: [
                {
                    question: "¿Cuáles son las comisiones para reservas de Hospedate?",
                    answer: "Para Reservas Directas, no hay ninguna comisión. Para Reservas Protegidas, Hospedate funciona con un sistema de comisión doble como hacen otras plataformas. Se calcula un subtotal con el precio de una propiedad por la cantidad de noches de la estadía. Sobre este subtotal, se aplica una comisión del 7% para el huésped. El anfitrión paga una comisión del 1.5% del subtotal."
                },
                {
                    question: "¿En qué se diferencian las comisiones de Hospedate versus otras plataformas?",
                    answer: "Con Hospedate tenés la posibilidad de no pagar comisiones. Eligiendo Reservas Directas, no hay comisión alguna. La otra diferencia, respecto de las Reservas Protegidas, si tomamos como referencia a Airbnb, las comisiones de Hospedate son 50% menores. Por cada reserva protegida que hacés con Hospedate, estás ahorrando un 50% sin importar si sos huésped o anfitrión."
                },
                {
                    question: "Como anfitrión, ¿De cuánto son los cargos por usar Hospedate?",
                    answer: "Como anfitrión en Hospedate, existen dos cargos. El primero es una suscripción como anfitrión de 3 USD por mes. El segundo, consiste en las comisiones en Reservas Protegidas."
                },
                {
                    question: "Como anfitrión, ¿De cuánto son los cargos por usar Hospedate únicamente con Reservas Directas?",
                    answer: "Si únicamente usás Hospedate con Reservas Directas, entonces el único cargo aplicable sería el de la suscripción como anfitrión de 3 USD por mes. Las Reservas Directas que hagas no tienen ninguna comisión."
                },
            ]
        },
        {
            title: "Beneficios en proceso de beta",
            questions: [
                {
                    question: "¿Cuáles son los beneficios por usar Hospedate durante el proceso de beta?",
                    answer: "Por usar Hospedate durante el proceso de beta, tenés múltiples beneficios. Si sos anfitrión, contás con una bonificación del 100% para usar Hospedate en modo anfitrión durante 12 meses. Además, por cada usuario que invites, podés ganar hasta USDT $4 como crédito disponible para usar en la plataforma! El crédito puede ser usado como vos quieras, ya sea para hacer tus comisiones aún más bajas o bien pagar parte de una reserva protegida."
                },
            ]
        }
    ]
}



const FaqPage = () => {

    const [user, setUser] = useContext(UserContext);
    const [faqInfo, setFaqInfo] = useState(initialFaqInfo)

    if (isGuest(user)) {
        Router.push('/')
    }

    const toggleDisplay = (clickedQuestion) => {
        setFaqInfo({
            ...faqInfo,
            sections: faqInfo.sections.map(section => ({
                ...section,
                questions: section.questions.map(question => ({
                    ...question,
                    visible: question.question === clickedQuestion.question ? (!question.visible) : question.visible,
                }))
            }))
        })
    }

    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true}>
                {isLoggedIn(user) && (
                    <ContentCard>
                        <Link className="text-hosp-light-blue hover:underline cursor-pointer mb-4 block" href="/help">← volver a la página principal de ayuda</Link>
                        <h1 className="text-3xl mb-4">Preguntas Frecuentes</h1>
                        {faqInfo.sections.map(section => (
                            <div key={section.title}>
                                <h2 className="mt-8">{section.title}</h2>
                                {section.questions.map(question => (
                                    <div key={question.question}>
                                        <h3><p className="pl-4 text-hosp-light-blue hover:underline cursor-pointer mb-2 block text-lg" onClick={(e) => {toggleDisplay(question)}}>{question.question}</p></h3>
                                        <p className={"pl-8 mt-2 mb-4 overflow-hidden transition-[max-height] delay-0 ease-linear " + (question.visible ? "max-h-[500px] duration-[3000ms]" : "max-h-0 duration-[500ms]")}>{question.answer}</p>
                                    </div>
                                ))}
                            </div>
                        ))}
                    </ContentCard>
                )}
            </SmartContainer>
        </div>
    )
}

export default FaqPage
