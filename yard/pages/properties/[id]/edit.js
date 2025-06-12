import NavBar from '../../../components/navBar'
import FormInputIcon from "../../../components/formInputIcon";
import FormTextArea from "../../../components/formTextArea";
import {useEffect, useState} from "react";
import Router, {useRouter} from "next/router";
import FormSelect from "../../../components/formSelect";
import utils from "../../../utils/utils";
import FormButton from "../../../components/formButton";
import componentsUtils from "../../../components/utils";
import PriceCalculator from "../../../components/priceCalculator";
import toast from "react-hot-toast";
import FormCheckbox from "../../../components/formCheckbox";
import SmartContainer from "../../../components/smartContainer";


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


async function handleSaveChanges(e, propertyId, changedFields, setSaveButtonEnabled) {
    e.preventDefault()
    setSaveButtonEnabled(false)
    if (Object.keys(changedFields).length === 0) {
        toast.error("no hay cambios")
        setSaveButtonEnabled(true)
        return
    }
    toast.loading("Guardando cambios", {
        id: "editingProperty"
    })
    const editProp = await fetch(
        `/api/properties/${propertyId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(changedFields)
        }
    )
    const response = await editProp.json()
    if (response.error) {
        toast.error("Error al guardar los cambios", {
            id: "editingProperty"
        })
        setSaveButtonEnabled(true)
    } else {
        toast.success("Cambios guardados", {
            id: "editingProperty"
        })
        setTimeout(() => {
            Router.push(`/properties/${propertyId}`)
        }, 500)
    }
}

function EditProperty() {

    const router = useRouter();
    const [prop, setProp] = useState({});
    const [saveButtonEnabled, setSaveButtonEnabled] = useState(false)

    const [title, setTitle] = useState("Cargando título")
    const [description, setDescription] = useState("Cargando descripción")
    const [maxGuests, setMaxGuests] = useState(0)
    const [price, setPrice] = useState(0)
    const [halfBathrooms, setHalfBathrooms] = useState("")
    const [bedrooms, setBedrooms] = useState("")

    const [accommodation, setAccommodation] = useState("");
    const [location, setLocation] = useState("")
    const [wifi, setWifi] = useState("")
    const [tv, setTv] = useState("")
    const [parking, setParking] = useState("")

    const [microwave, setMicrowave] = useState(false)
    const [oven, setOven] = useState(false)
    const [kettle, setKettle] = useState(false)
    const [toaster, setToaster] = useState(false)
    const [coffeeMachine, setCoffeeMachine] = useState(false)
    const [airConditioning, setAirConditioning] = useState(false)
    const [heating, setHeating] = useState(false)
    const [pool, setPool] = useState(false)
    const [gym, setGym] = useState(false)

    useEffect(() => {
        const changedFields = collectChangedFields()
        setSaveButtonEnabled(Boolean(Object.keys(changedFields).length > 0))
    }, [
        title,
        description,
        maxGuests,
        price,
        halfBathrooms,
        bedrooms,
        accommodation,
        location,
        wifi,
        tv,
        parking,
        microwave,
        oven,
        kettle,
        toaster,
        coffeeMachine,
        airConditioning,
        heating,
        pool,
        gym,
    ])

    function collectChangedFields() {
        const candidates = {
            "title": title,
            "description": description,
            "maxGuests": maxGuests,
            "price": price,
            "accommodation": accommodation,
            "location": location,
            "wifi": wifi,
            "tv": tv,
            "parking": parking,
            "microwave": microwave,
            "oven": oven,
            "kettle": kettle,
            "toaster": toaster,
            "coffeeMachine": coffeeMachine,
            "airConditioning": airConditioning,
            "heating": heating,
            "pool": pool,
            "gym": gym,
            "halfBathrooms": halfBathrooms,
            "bedrooms": bedrooms,
        }
        const changedFields = {}
        Object.keys(candidates).map((field) => {
            if (typeof candidates[field] === 'string' && typeof prop[field] === 'string' && candidates[field].trim() !== prop[field].trim()) {
                changedFields[field] = candidates[field].trim()
            } else if (typeof prop[field] === 'number' && parseInt(candidates[field], 10) !== prop[field]) {
                changedFields[field] = parseInt(candidates[field], 10)
            } else if (prop[field] === null && typeof candidates[field] === 'number') {
                changedFields[field] = candidates[field]
            } else if (prop[field] === null && typeof candidates[field] === 'string' && candidates[field] !== '') {
                changedFields[field] = candidates[field]
            } else if (typeof prop[field] === 'string' && typeof candidates[field] === 'boolean' && prop[field] === 'not_available' && candidates[field]) {
                changedFields[field] = "available"
            } else if (typeof prop[field] === 'string' && typeof candidates[field] === 'boolean' && prop[field] === 'available' && !candidates[field]) {
                changedFields[field] = "not_available"
            } else if (prop[field] === null && typeof candidates[field] === 'boolean' && candidates[field]) {
                changedFields[field] = "available"
            }
        })
        return changedFields
    }

    function setEditForm(property) {
        setTitle(property.title)
        setDescription(property.description)
        setMaxGuests(property.maxGuests)
        setPrice(property.price)
        setHalfBathrooms(property.halfBathrooms === null ? "" : property.halfBathrooms)
        setBedrooms(property.bedrooms === null ? "" : property.bedrooms)
        // amenities
        setAccommodation(property.accommodation || "")
        setLocation(property.location || "")
        setWifi(property.wifi || "")
        setTv(property.tv || "")
        setParking(property.parking || "")
        // boolean attributes
        setMicrowave(property.microwave === "available")
        setOven(property.oven === 'available')
        setKettle(property.kettle === 'available')
        setToaster(property.toaster === 'available')
        setCoffeeMachine(property.coffeeMachine === 'available')
        setAirConditioning(property.airConditioning === 'available')
        setHeating(property.heating === 'available')
        setPool(property.pool === 'available')
        setGym(property.gym === 'available')
    }

    useEffect(() => {
        (async () => {
            if (router.query.id) {
                const prop = await fetchProperty(router.query.id);
                setProp(prop);
                setEditForm(prop)
            }
        })()
    }, [router.isReady])

    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true}>
                <h1 className="text-2xl md:mb-7">Editar propiedad</h1>
                <form onSubmit={(e) => handleSaveChanges(e, prop.id, collectChangedFields(), setSaveButtonEnabled)}>
                    <h2>Título y descripción</h2>
                    <FormInputIcon icon="title" name="title" id="title" widthClasses="w-full" value={title} onChange={(e) => {setTitle(e.target.value)}} required />
                    <FormTextArea icon="description" widthClasses="w-full" extraClasses="h-[300px]" name="description" id="description" value={description} onChange={(e) => {setDescription(e.target.value)}} required />

                    <h2 className="mt-5 mb-2">Capacidad y precio</h2>
                    <div className="flex flex-col tablet:flex-row flex-wrap">
                        <FormInputIcon icon="group" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" type="number" name="max_guests" id="max_guests" value={maxGuests} onChange={(e) => {setMaxGuests(e.target.value)}} required />
                        <FormInputIcon icon="attach_money" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" type="number" name="price" id="price" value={price} onChange={(e) => {setPrice(e.target.value)}} placeholder="Precio por noche (USDT)" required />
                        <FormSelect icon="wc" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="halfBathrooms" id="form-halfBathrooms" value={halfBathrooms} onChange={e => setHalfBathrooms(parseInt(e.target.value, 10))}>`
                            {prop.halfBathrooms === null ? (
                                <option value="">Cantidad de baños</option>
                            ) : null}
                            {[...Array(21).keys()].map(hb => (
                                <option key={hb} value={hb}>{hb / 2}</option>
                            ))}
                        </FormSelect>
                        <FormSelect icon="bed" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="bedrooms" id="form-bedrooms" value={bedrooms} onChange={e => setBedrooms(parseInt(e.target.value, 10))}>`
                            {prop.bedrooms === null ? (
                                <option value="">Dormitorios</option>
                            ) : null}
                            {[...Array(11).keys()].map(br => (
                                <option key={br} value={br}>{br}</option>
                            ))}
                        </FormSelect>
                    </div>
                    <div className="flex justify-center">
                        <PriceCalculator price={price} extraClasses="md:w-[600px]" />
                    </div>

                    <h2 className="mt-5 mb-2">Servicios</h2>
                    <div className="grid grid-cols-1 tablet:grid-cols-2 desktop:grid-cols-3">
                        <FormSelect icon="home" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="accommodation" id="form-accommodation" value={accommodation} onChange={e => setAccommodation(e.target.value)}>`
                            {prop.accommodation === null ? (
                                <option value="">Tipo de alojamiento</option>
                            ) : null}
                            {utils.accommodations.map(acc => (
                                <option key={acc} value={acc}>{componentsUtils.humanizeAccommodation(acc)}</option>
                            ))}
                        </FormSelect>
                        <FormSelect icon="map" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="location" id="form-location" value={location} onChange={e => setLocation(e.target.value)}>`
                            {prop.location === null ? (
                                <option value="">Ubicación</option>
                            ) : null}
                            {utils.locations.map(loc => (
                                <option key={loc} value={loc}>{componentsUtils.humanizeLocation(loc)}</option>
                            ))}
                        </FormSelect>
                        <FormSelect icon="wifi" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="wifi" id="form-wifi" value={wifi} onChange={e => setWifi(e.target.value)}>`
                            {prop.wifi === null ? (
                                <option value="">Opción de Wifi</option>
                            ) : null}
                            {utils.wifiOptions.map(wifiOpt => (
                                <option key={wifiOpt} value={wifiOpt}>{componentsUtils.humanizeWifiOption(wifiOpt)}</option>
                            ))}
                        </FormSelect>
                        <FormSelect icon="tv" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="wifi" id="form-tv" value={tv} onChange={e => setTv(e.target.value)}>`
                            {prop.tv === null ? (
                                <option value="">Opción de Tv</option>
                            ) : null}
                            {utils.tvOptions.map(tvOpt => (
                                <option key={tvOpt} value={tvOpt}>{componentsUtils.humanizeTvOption(tvOpt)}</option>
                            ))}
                        </FormSelect>
                        <FormSelect icon="garage_home" containerExtraClasses="md:mr-4" widthClasses="min-w-[250px]" name="parking" id="form-parking" value={parking} onChange={e => setParking(e.target.value)}>`
                            {prop.parking === null ? (
                                <option value="">Estacionamiento</option>
                            ) : null}
                            {utils.parkingOptions.map(parkingOpt => (
                                <option key={parkingOpt} value={parkingOpt}>{componentsUtils.humanizeParkingOption(parkingOpt)}</option>
                            ))}
                        </FormSelect>
                    </div>
                    <div className="mt-5 grid grid-cols-1 tablet:grid-cols-2 desktop:grid-cols-3">
                        <FormCheckbox label="Microondas" icon="microwave" checked={microwave} onChange={(e) => setMicrowave(e.target.checked)} />
                        <FormCheckbox label="Horno" icon="oven_gen" checked={oven} onChange={(e) => setOven(e.target.checked)} />
                        <FormCheckbox label="Pava eléctrica" icon="kettle" checked={kettle} onChange={(e) => setKettle(e.target.checked)} />
                        <FormCheckbox label="Tostadora" icon="breakfast_dining" checked={toaster} onChange={(e) => setToaster(e.target.checked)} />
                        <FormCheckbox label="Máquina de café" icon="coffee_maker" checked={coffeeMachine} onChange={(e) => setCoffeeMachine(e.target.checked)} />
                        <FormCheckbox label="Aire acondicionado" icon="mode_fan" checked={airConditioning} onChange={(e) => setAirConditioning(e.target.checked)} />
                        <FormCheckbox label="Calefacción" icon="fireplace" checked={heating} onChange={(e) => setHeating(e.target.checked)} />
                        <FormCheckbox label="Piscina" icon="pool" checked={pool} onChange={(e) => setPool(e.target.checked)} />
                        <FormCheckbox label="Gimnasio" icon="exercise" checked={gym} onChange={(e) => setGym(e.target.checked)} />
                    </div>

                    <div className="flex justify-center">
                        <FormButton disabled={!saveButtonEnabled} type="submit" text="Guardar cambios" />
                    </div>
                </form>
            </SmartContainer>
        </div>
    );
}

export default EditProperty