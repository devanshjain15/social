import { useParams } from "react-router";
import { BASE_API } from "../main";

export const ConfirmUser = () => {
    const {token = ''} = useParams(); 
    const handleClick = async () => { 
        try {
            const res = await fetch(`${BASE_API}/users/activate/${token}`, { 
                method: "PUT"
            }); 

            if (res.ok){ 
                console.log("✅ Success!"); 
            } else { 
                console.log("❌ Request Failed!"); 
            }
        } catch (error) {
            console.log(error)
        }
    }
    return (
        <>
        <div>Confirm User</div>
        <button onClick={handleClick}>Click here to confirm</button>
        </>
    )
}
