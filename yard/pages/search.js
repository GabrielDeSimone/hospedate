import NavBar from '../components/navBar';
import { SearchParamsContext } from '../components/searchParamsProvider';
import utils from '../components/utils'
import { useState, useEffect, useContext } from 'react';
import { useRouter } from 'next/router';
import PropertyGrid from '../components/propertyGrid';

async function fetchProperties(searchParams) {
    const city = searchParams.city;
    const checkinDate = searchParams.checkinDate;
    const checkoutDate = searchParams.checkoutDate;
    const guests = searchParams.guests;

    const getProperties = await fetch(
        "/api/properties/search?city=" +
        city +
        "&checkinDate=" +
        checkinDate +
        "&checkoutDate=" +
        checkoutDate +
        "&guests=" +
        guests
    )
    const response = await getProperties.json()

    if (response.data) {
        return response.data
    } else {
        return []
    }
}

function getSearchParams() {
    const search = window.location.search;
    const params = new URLSearchParams(search);

    return {
        city: params.get('city'),
        checkinDate: params.get('checkinDate'),
        checkoutDate: params.get('checkoutDate'),
        stayDays: utils.getDaysDifference(
            params.get('checkinDate'),
            params.get('checkoutDate')
        ),
        guests: params.get('guests'),
    }
}

const Search = (props) => {
    const [properties, setProperties] = useState([]);
    const [searchParams, setSearchParams] = useContext(SearchParamsContext);
    const [stayParams, setStayParams] = useState({});
    const router = useRouter();

    useEffect(() => {
        const urlParams = getSearchParams();
        setSearchParams(urlParams)
        setStayParams({
            checkinDate: urlParams.checkinDate,
            checkoutDate: urlParams.checkoutDate,
            guests: urlParams.guests,
            stayDays: urlParams.stayDays,
        })

        async function asyncEffect() {
            const propertiesFetched = await fetchProperties(urlParams);
            setProperties(propertiesFetched);
        }
        asyncEffect();
    }, [router.query]);

    return (
        <div>
            <NavBar showSearchParams={true} />
            {properties.length > 0 ? (
                <PropertyGrid properties={properties} stayParams={stayParams} />
            ) : (
                <p className="text-center mt-6 text-gray-500 text-lg">No hay resultados</p>
            )}

        </div>
    )
}

export default Search
