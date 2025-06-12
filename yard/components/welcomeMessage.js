import NavBar from "./navBar";
import ContentCard from "./contentCard";
import Link from "next/link";
const WelcomeMessage = () => {

    return (
        <div>
            <ContentCard extraClasses="p-6 mt-6 mx-6">
                <h1 className="text-3xl mb-4">¡Te damos la bienvenida a Hospedate! 🙌</h1>
                <p className="leading-8">Estamos muy emocionados de que seas parte de nuestro proceso de beta. Para nosotros es muy importante el uso que le des a nuestro producto, ya que tu feedback será clave y transformará el futuro de Hospedate.</p>
                <p className="leading-8">Nuestro equipo ha trabajado mucho para poder brindarte un producto que creemos que es parte de un nuevo paradigma en la forma en que la gente usa su dinero, con mayor libertad financiera y autonomía individual.</p>
                <h2 className="mt-7 mb-3">Qué esperar de este proceso de beta</h2>
                <ul className="list-disc px-10">
                    <li className="leading-8"><span className="font-bold">✅ Acceso exclusivo</span>: Como usuario pionero, estarás entre las primeras personas en probar las nuevas funcionalidades antes de que nuestro producto esté disponible a todo público.</li>
                    <li className="leading-8"><span className="font-bold">🤝 Colaboración</span>: Te invitamos a que participes en nuestra comunidad compartiendo tus opiniones o reportando cualquier inconveniente con el que te encuentres. Podés comunicarte con nosotros por X (ex Twitter) en <Link className="text-hosp-light-blue hover:underline cursor-pointer" target="_blank" href="https://twitter.com/hospedate_app">@hospedate_app</Link>, o al correo electrónico y teléfonos descritos en la sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help">Ayuda</Link>.</li>
                    <li className="leading-8"><span className="font-bold">🌱 Mejora continua</span>: Estamos comprometidos a mejorar Hospedate según tu feedback. Conocer y mejorar tu experiencia es nuestra prioridad número 1. Por eso, podés contactarnos no sólo para resolver un inconveniente puntual, sino para proponer ideas sobre cómo podemos mejorar nuestro producto de la forma que mejor te satisfaga.</li>
                    <li className="leading-8"><span className="font-bold">🚀 Beneficios exclusivos</span>: Por entrar durante nuestro proceso de beta, podés disfrutar de Hospedate en modo anfitrión por <span className="font-bold">12 meses sin cargo</span>.</li>
                    <li className="leading-8"><span className="font-bold">💸 Además</span>: Ganá hasta <span className="font-bold">4 USDT</span> por cada usuario que invites como crédito para usar en Hospedate 😎.</li>
                </ul>
                <h2 className="mt-7 mb-3">¿Preguntas?</h2>
                <p>No dudes en visitar nuestra sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">Preguntas frecuentes</Link> para saber más acerca de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">Cómo funcionan los beneficios</Link>, <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">cuáles son los precios de Hospedate</Link>, <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">cómo se diferencia Hospedate de otras plataformas</Link>, entre otras preguntas.</p>
                <h2 className="mt-7 mb-3">Estamos a un click de distancia</h2>
                <p>Queremos que sepas que estamos disponibles para asistirte con cualquier inquietud o inconveniente que puedas tener. Podés encontrar los datos de contacto de nuestro equipo de soporte en la página de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help">Ayuda</Link>.</p>
            </ContentCard>
        </div>
    )
}

export default WelcomeMessage
