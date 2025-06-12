import ContentBox from '../components/contentBox';
import NavBar from '../components/navBar';
import SearchParamsWindow from '../components/searchParamsWindow';
import TemporaryIndex from '../components/temporaryIndex';
import { UserContext, isLoggedIn } from '../components/userProvider';
import { useContext } from 'react';
import WelcomeHosts from "../components/welcomeHosts";
import SmartContainer from "../components/smartContainer";
import WelcomeMessage from "../components/welcomeMessage";

const IndexPage = () => {
    const [user, setUser] = useContext(UserContext);

    return (
        <div>
            {
                isLoggedIn(user) ? (
                    <div>
                        <NavBar />
                        <SmartContainer yAxisSpaced={true}>
                            <WelcomeMessage />
                        </SmartContainer>
                    </div>
                ) : (
                    <WelcomeHosts />
                )
            }
        </div>
    )
}

export default IndexPage
