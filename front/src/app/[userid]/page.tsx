export default function Page({ params }: { params: { userid: string } }) {
  return <div>My Post: {params.userid}</div>;
}
