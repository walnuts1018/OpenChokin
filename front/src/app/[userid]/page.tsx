import { MypageComponent } from "../mypage/MypageComponent";
export default function Page({ params }: { params: { userid: string } }) {
  return <MypageComponent userID={params.userid} />;
}
