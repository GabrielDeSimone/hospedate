import { useContext, useEffect, useState } from "react"
import FormButton from "./formButton"
import FormInputIcon from "./formInputIcon"
import Router from 'next/router'
import { SearchParamsContext } from "./searchParamsProvider"
import utils from "../utils/utils";
import FormSelect from "./formSelect";
import compUtils from "./utils"
import toast from "react-hot-toast";

function handleSubmit(callback) {
    function innerHandleSubmit(event) {
        event.preventDefault()
        // collect all form data
        const data = {
            city: event.target.city.value,
            checkinDate: event.target.checkinDate.value,
            checkoutDate: event.target.checkoutDate.value,
            guests: parseInt(event.target.guests.value, 10) || 2,
        }

        if (compUtils.getDaysDifference(compUtils.getLocalTodayDate(), data.checkinDate) < 0) {
            toast.error("La fecha de check-in no puede estar en el pasado")
            return
        }

        if (compUtils.getDaysDifference(data.checkinDate, data.checkoutDate) < 1) {
            toast.error("La fecha de check-out debe ser posterior a la de check-in")
            return
        }


        // redirect to /search with all data as query params
        Router.push({
            pathname: `/search`,
            query: data,
        })
    }

    return (event) => {
        innerHandleSubmit(event)
        if (callback) {
            callback()
        }
    }
}

const SearchParamsWindow = (props) => {
    const [searchParams, setSearchParams] = useContext(SearchParamsContext);
    const [city, setCity] = useState("")
    const [checkinDate, setCheckinDate] = useState("")
    const [checkoutDate, setCheckoutDate] = useState("")
    const [guests, setGuests] = useState("")

    useEffect(() => {
        if (searchParams) {
            setCity(searchParams.city)
            setCheckinDate(searchParams.checkinDate)
            setCheckoutDate(searchParams.checkoutDate)
            setGuests(searchParams.guests)
        }
    }, []);

    return (
        <div>
            <form className="flex flex-col items-center" onSubmit={handleSubmit(props.onSubmit)}>
                <FormSelect icon="location_on" name="city" id="form-city" value={city} onChange={e => setCity(e.target.value)} required>
                    <option disabled="disabled" value="">Ciudad</option>
                    {utils.cities.map(city => (
                        <option key={city} value={city}>{city}</option>
                    ))}
                </FormSelect>
                <FormInputIcon
                    type="date"
                    icon="date_range"
                    placeholder="Fecha de check-in"
                    name="checkinDate"
                    id="form-date-in"
                    value={checkinDate}
                    onChange={e => setCheckinDate(e.target.value)}
                    required={true} />
                <FormInputIcon
                    type="date"
                    icon="date_range"
                    placeholder="Fecha de check-out"
                    name="checkoutDate"
                    id="form-date-out"
                    value={checkoutDate}
                    onChange={e => setCheckoutDate(e.target.value)}
                    required={true} />
                <FormInputIcon
                    type="number"
                    icon="group"
                    placeholder="Huéspedes"
                    name="guests"
                    value={guests}
                    onChange={e => setGuests(e.target.value)}
                    id="form-guests" />
                <div className="mt-8 w-[80%] sm:w-auto">
                    <FormButton type="submit" text="Buscar" />
                </div>
            </form>
        </div>
    )
}

export default SearchParamsWindow