import NavBar from '../../components/navBar'
import {useContext, useEffect, useState} from 'react'
import {isGuest, UserContext} from '../../components/userProvider'
import Router, { useRouter } from "next/router";
import Slider from "../../components/slider";
import PropertyInfoBar from "../../components/property/infoBar";
import PriceDisplayer from "../../components/property/pricedisplayer";
import componentsUtils from "../../components/utils";
import utils from '../../utils/utils'
import FormButton from "../../components/formButton";
import toast from "react-hot-toast";
import SmartContainer from "../../components/smartContainer";
import ContentCard from "../../components/contentCard";

const apiErrors = utils.apiErrors

async function fetchProperty(propId) {
    const getProperty = await fetch(
        `/api/properties/${propId}`
    )
    const response = await getProperty.json()

    if (response.data) {
        return response.data
    } else if (getProperty.status === 404) {
        Router.push('/404')
    } else {
        return {}
    }
}

async function handleEditProp(e, propertyId) {
    e.preventDefault()
    Router.push(`/properties/${propertyId}/edit`)
}

async function handleRemoveProp(e, propertyId) {
    e.preventDefault()
    toast.loading("Eliminando propiedad", {
        id: "removingProperty"
    })
    const removeProp = await fetch(
        `/api/properties/${propertyId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
        }
    )
    const response = await removeProp.json()
    if (response.error && response.error === apiErrors.PropertyActiveOrders) {
        toast.error("La propiedad tiene reservas activas", {
            id: "removingProperty"
        })
    } else if (response.error) {
        toast.error("Hubo un problema al eliminar la propiedad", {
            id: "removingProperty"
        })
    } else {
        toast.success("La propiedad se eliminó exitosamente", {
            id: "removingProperty"
        })
        setTimeout(() => {
            Router.push(`/hosts`)
        }, 1000)
    }
}

function getStayParams() {
    const search = window.location.search;
    const params = new URLSearchParams(search);

    const checkinDate = params.get('checkinDate') || componentsUtils.getTodayDate();
    const checkoutDate = params.get('checkoutDate') || getTwoDaysAfterDate();

    return {
        checkinDate,
        checkoutDate,
        stayDays: componentsUtils.getDaysDifference(checkinDate, checkoutDate),
        guests: parseInt(params.get('guests'), 10) || 2,
    }
}

function getTwoDaysAfterDate() {
    return new Date(Date.now() + 2 * 24 * 60 * 60 * 1000).toISOString().slice(0, 10)
}


function Property() {
    const [user, setUser] = useContext(UserContext);
    const [prop, setProp] = useState({});
    const router = useRouter();
    const [isOwnerSeeing, setIsOwnerSeeing] = useState(false)
    const [stayParams, setStayParams] = useState({
        checkinDate: componentsUtils.getTodayDate(),
        checkoutDate: getTwoDaysAfterDate(),
        stayDays: componentsUtils.getDaysDifference(componentsUtils.getTodayDate(), getTwoDaysAfterDate()),
        guests: 2
    });

    useEffect(() => {
        const stayParams = getStayParams();
        setStayParams(stayParams)
    }, [router.query])

    useEffect(() => {
        (async () => {
            if (router.query.id) {
                const prop = await fetchProperty(router.query.id);
                if (prop && prop.images) {
                    prop.images = prop.images.map(img => img+"?im_w=1200")
                }
                setProp(prop);
                setIsOwnerSeeing(Boolean(prop.user_id && user && prop.user_id === user.id))
            }
        })()
    }, [router.isReady, user])

    useEffect(() => {
        if (isGuest(user)) {
            Router.push('/')
        }
    }, [user])

    return (
        <div>
            <NavBar />
            <SmartContainer noTabletWidthLimit={true}>
                {isOwnerSeeing && (
                    <div className="flex p-6">
                        <form className="mr-4" onSubmit={(e) => {handleEditProp(e, prop.id)}}>
                            <FormButton type="submit" text="Editar propiedad" extraClasses="sm:w-[300px] ml-auto mr-auto block" />
                        </form>
                        <form onSubmit={(e) => {handleRemoveProp(e, prop.id)}}>
                            <FormButton type="submit" text="Eliminar propiedad" extraClasses="sm:w-[300px] ml-auto mr-auto block" />
                        </form>
                    </div>
                )}
                <h1 className="text-2xl mt-5 mb-5">{prop.title}</h1>
                <div className="mb-4">
                    <Slider images={prop.images}/>
                </div>
                <div className="flex flex-col tablet:flex-row">
                    <div id="propertyInformation" className="tablet:w-2/3 laptop:w-3/4">
                        <div className="px-3 py-1">
                            <PropertyInfoBar
                                property={prop}
                                checkinDate={stayParams.checkinDate}
                                checkoutDate={stayParams.checkoutDate}
                                isOwner={Boolean(prop.user_id && prop.user_id === user.id)}
                            />
                        </div>
                        <div>
                            {prop.description ? (
                                <div className="leading-8 text-gray-900 mt-8 p-6">{(
                                    prop.description.split('\n').map((descParagraph, index) => (
                                        (descParagraph === '' ? <br key={index} /> : <p key={index}>{descParagraph}</p>)
                                    ))
                                )}</div>
                            ) : null}
                        </div>
                    </div>
                    <div className="tablet:w-1/3 laptop:w-1/4">
                        {!isOwnerSeeing && (
                            <ContentCard shadow={true}>
                                <PriceDisplayer stayParams={stayParams} price={prop.price} propertyId={prop.id} />
                            </ContentCard>
                        )}
                    </div>
                </div>
            </SmartContainer>
        </div>
    );
}

export default Property