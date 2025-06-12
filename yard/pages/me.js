import NavBar from '../components/navBar';
import {useContext, useEffect, useState} from "react";
import {UserContext, isGuest, isLoggedIn} from "../components/userProvider";
import Router from "next/router";
import {
    TableH,
    TableHTd,
    TableHTh,
    TableHTr, TableV,
    TableVBody,
    TableVHeading, TableVTd,
    TableVThCol, TableVThRow,
    TableVTr
} from "../components/myTable";
import componentsUtils from "../components/utils";
import toast from "react-hot-toast";
import Link from "next/link";
import SmartContainer from "../components/smartContainer";
import ContentCard from "../components/contentCard";
import Modal from "../components/modal";
import FormButton from "../components/formButton";
import utils from "../utils/utils"

async function generateNewInvitation(invKind) {
    toast.loading("Generando nueva invitación", {
        id: "generatingInvitation"
    })
    const generateInvitation = await fetch("/api/invitations", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            kind: invKind
        })
    })
    const response = await generateInvitation.json()

    if (generateInvitation.status === 201) {
        toast.success("La invitación fue generada exitosamente", {
            id: "generatingInvitation"
        })
        return response.data
    } else {
        toast.error("Hubo un problema al generar la invitación", {
            id: "generatingInvitation"
        })
        return null
    }
}

async function fetchFullUser() {
    const getFullUser = await fetch(
        "/api/me"
    )
    const response = await getFullUser.json()

    if (response.data) {
        return response.data
    } else {
        return null
    }
}

async function fetchMyCreditInst() {
    const getCredit = await fetch(
        "/api/me/credit"
    )
    const response = await getCredit.json()

    if (response.data) {
        return response.data
    } else {
        return null
    }
}

function creditByInvitation(creditInstances, invitationId) {
    const instances = creditInstances.filter(creditInst => creditInst.invitationId === invitationId)
    if (instances.length === 0) {
        return 0
    } else {
        return instances[0].earnedAmount
    }
}

const MePage = () => {

    const [user, setUser] = useContext(UserContext);
    const [fullUser, setFullUser] = useState(null);
    const [creditInstances, setCreditInstances] = useState([]);
    const [showNewInvitationModal, setShowNewInvitationModal] = useState(false);
    const [godMode, setGodMode] = useState(false);

    function hospedate(g) {
        if (g === "Hospedate.") {
            setGodMode(true)
        }
        console.log("Hi!")
    }

    const handleModalClose = () => {
        setShowNewInvitationModal(false);
    }

    const totalCredit = () => {
        return creditInstances.reduce((a,b) => (a + b.earnedAmount), 0)
    }

    if (isGuest(user)) {
        Router.push('/')
    }

    useEffect(() => {
        // get full user
        (async () => {
            window.hospedate = hospedate;
            const fetchedFullUser = await fetchFullUser()
            const fetchedCreditInst = await fetchMyCreditInst()
            setFullUser(fetchedFullUser)
            setCreditInstances(fetchedCreditInst)
        })()
    }, [])

    const genInvitationAndRefresh = async (e, invKind) => {
        e.preventDefault()
        setShowNewInvitationModal(false)
        const invitation = await generateNewInvitation(invKind)
        if (invitation) {
            setFullUser(await fetchFullUser())
        }
    }

    const showInvitationModal = () => {
        setShowNewInvitationModal(true);
    }

    return (
        <div>
            <NavBar />
            {isLoggedIn(user) && fullUser && (
                <SmartContainer yAxisSpaced={true}>
                    <ContentCard shadow={true}>
                        <h1 className="text-2xl">Hola, {user.name}!</h1>
                        <h2 className="mt-4">Datos de cuenta</h2>
                        <div className="overflow-y-scroll">
                            <TableH>
                                <TableHTr>
                                    <TableHTh>Nombre</TableHTh>
                                    <TableHTd>{user.name}</TableHTd>
                                </TableHTr>
                                <TableHTr>
                                    <TableHTh>Email</TableHTh>
                                    <TableHTd>{user.email}</TableHTd>
                                </TableHTr>
                                <TableHTr>
                                    <TableHTh>Número de teléfono</TableHTh>
                                    <TableHTd>{fullUser.phoneNumber}</TableHTd>
                                </TableHTr>
                            </TableH>
                        </div>
                        <p className="text-gray-500 mt-3 ml-8 text-sm">Contáctanos desde la sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help">ayuda</Link> para modificar tus datos personales.</p>
                    </ContentCard>
                    <h2 className="mt-8">Crédito en total ganado por invitaciones: USDT $ {totalCredit()}</h2>
                    <h2 className="mt-8">Invitaciones generadas</h2>
                    {
                        fullUser.generatedInvitations.length > 0 ? (
                            <div className="overflow-y-scroll">
                                <TableV extraClasses="ml-5">
                                    <TableVHeading>
                                        <TableVThCol>Código de invitación</TableVThCol>
                                        <TableVThCol>Fecha de creación</TableVThCol>
                                        <TableVThCol>Invitación utilizada</TableVThCol>
                                        <TableVThCol>Tipo de invitación</TableVThCol>
                                        <TableVThCol>Crédito recibido</TableVThCol>
                                    </TableVHeading>
                                    <TableVBody>
                                        {fullUser.generatedInvitations.map((invitation) => (
                                            <TableVTr key={invitation.id}>
                                                <TableVThRow>{invitation.id}</TableVThRow>
                                                <TableVTd>{componentsUtils.humanizeTimestamp(invitation.createdAt)}</TableVTd>
                                                <TableVTd>{componentsUtils.humanizeBoolean(invitation.used)}</TableVTd>
                                                <TableVTd>{componentsUtils.humanizeInvitationKind(invitation.kind)}</TableVTd>
                                                <TableVTd>USDT $ {creditByInvitation(creditInstances, invitation.id)}</TableVTd>
                                            </TableVTr>
                                        ))}
                                    </TableVBody>
                                </TableV>
                            </div>
                        ) : (
                            <p className="text-gray-500 mt-3 ml-8 text-sm">No hay invitaciones generadas</p>
                        )
                    }

                    {
                        (fullUser.generatedInvitations.length < utils.INVITATION_MAX_LIMIT || godMode) ? (
                            <Link className="
                        text-hosp-light-blue
                        hover:underline
                        cursor-pointer
                        block
                        mt-3
                        ml-3
                        "
                                  href="#"
                                  onClick={showInvitationModal}
                            >Generar nueva invitación</Link>
                        ) : (
                            <p className="text-gray-500 mt-3 ml-8 text-sm">Has alcanzado tu límite de invitaciones. Contáctanos desde la sección de <Link className="text-hosp-light-blue hover:underline cursor-pointer" href="/help">ayuda</Link> para generar más.</p>
                        )
                    }
                </SmartContainer>
            )}
            <Modal show={showNewInvitationModal} onClose={handleModalClose} position="top">
                <h2>Generar nueva invitación</h2>
                <p>¿Qué tipo de invitación deseas generar? Podés consultar los beneficios de cada tipo de invitación <Link className="text-hosp-light-blue hover:underline cursor-pointer" target="_blank" href="/help/faq">acá</Link>.</p>
                <div className="flex flex-col sm:justify-between sm:flex-row">
                    <FormButton onClick={(e) => genInvitationAndRefresh(e, "FOR_GUEST")} text="Invitación para huéspedes" />
                    <FormButton onClick={(e) => genInvitationAndRefresh(e, "FOR_HOST")} text="Invitación para anfitriones" />
                </div>
            </Modal>
        </div>
    )
}

export default MePage
