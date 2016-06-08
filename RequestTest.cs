using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;

namespace Example.Http
{
    public class RequestTest
    {
        static bool requestFinished = false;
        static String requestUrl = null;
        static int intervalInMs = 100;

        static void Main(string[] args)
        {
            if (args.Length < 2)
            {
                Console.WriteLine("请输入URL和时间间隔");
                return;
            }

            requestUrl = args[0];
            intervalInMs = Convert.ToInt32(args[1]);

            Thread requestThread = new Thread(new ThreadStart(RequestUrl));

            Console.WriteLine("按回车开始");
            Console.ReadLine();

            requestThread.Start();

            Console.WriteLine("按回车停止");
            Console.ReadLine();

            requestFinished = true;
            requestThread.Join();

            Console.WriteLine("按回车退出");
            Console.ReadLine();
        }

        protected static void RequestUrl()
        {
            int requstCount = 0;
            int failedCount = 0;
            const string SUCCESS_RESPONSE = "\"code\":0";
            Stopwatch watch = new Stopwatch();
            
            DateTime startTime = DateTime.Now;
            while (!requestFinished)
            {
                watch.Reset();
                watch.Start();
                requstCount++;

                //GetHTMLTCP("http://61.147.186.54:8482/api/push?data=eyJtc2dDbGllbnRJZCI6IjU3M2FjZTJhYTJhOGE5MjFlYTc3ZjkyYSIsIm1zZ1RpbWUiOjE0NjM0NzE2NTg1NTksImRhdGFUeXBlIjoidV9saXZlX2RhdGEiLCJwYXJ0aWFsT3JkZXIiOjE0NjM0NzE2NTg1NTksInR5cGUiOjIsImxpZCI6IjU3M2FjYzkyOTdlNDY5MmZiNmZkOWI4NCIsInVpZCI6OTUwMDI1ODg2LCJ1c2VybmFtZSI6ImJhaWppMjIiLCJuaWNrIjoieGlhby3pmL/kuLkiLCJoZWFkZXJVcmwiOiJodHRwOi8vb3VydGltZXNwaWN0dXJlLmJzMmRsLnl5LmNvbS9vcGVyYXRpb25faGVhZGVyXzk1MDAwMDA3Ml8yX3Rlc3RfMTQ1NTU5NDE0MTIxOC5qcGciLCJzZXgiOjIsInZlcmlmaWVkIjpmYWxzZSwiZ3Vlc3RDb3VudCI6NiwidG90YWxHdWVzdENvdW50IjozLCJndWVzdFJhdGUiOjB9&topic=/testluoyanjie&pushAll=true");
                //GetHTMLTCP("http://test.mlc.yy.com/api/bcproxy_svr/broadcast?appId=100001&sign=&data={\"topic\":\"yy1\", \"message\":\"jsonssdssssssssssssssssssssssssssssssssffffffsdfffffffffffffffffffffffffffffftring56765675656777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777777\", \"alive_time\":12, \"is_order\":true}");
                String returnVal = GetHTMLTCP(requestUrl);
                if (!returnVal.Contains(SUCCESS_RESPONSE))
                {
                    Console.WriteLine("Failed: " + returnVal);
                    failedCount++;
                }

                Console.WriteLine(requstCount.ToString());
                //Console.WriteLine(returnVal);

                long timeInMs = watch.ElapsedMilliseconds;
                if (timeInMs < intervalInMs)
                {
                    Thread.Sleep(intervalInMs - (int)timeInMs);
                }
                watch.Stop();
            }
            TimeSpan elapesdTime = DateTime.Now.Subtract(startTime);
            
            Console.WriteLine("总请求数:" + requstCount);
            Console.WriteLine("失败请求数:" + failedCount);
            Console.WriteLine("总时间:" + elapesdTime.TotalSeconds);
            Console.WriteLine("平均每秒请求数:" + requstCount / elapesdTime.TotalSeconds);
        }

        private static string GetHTMLTCP(string URL)
        {
            string strHTML = "";
            TcpClient clientSocket = new TcpClient();
            Uri URI = new Uri(URL);
            clientSocket.Connect(URI.Host, URI.Port);
            StringBuilder RequestHeaders = new StringBuilder();
            RequestHeaders.AppendFormat("{0} {1} HTTP/1.1\r\n", "GET", URI.PathAndQuery);
            RequestHeaders.AppendFormat("Connection:close\r\n");
            RequestHeaders.AppendFormat("Host:{0}\r\n", URI.Host);
            RequestHeaders.AppendFormat("Accept:*/*\r\n");
            RequestHeaders.AppendFormat("Accept-Language:zh-cn\r\n");
            RequestHeaders.AppendFormat("User-Agent:Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.2; SV1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)\r\n\r\n");
            Encoding encoding = Encoding.Default;
            byte[] request = encoding.GetBytes(RequestHeaders.ToString());
            clientSocket.Client.Send(request);

            Stream readStream = clientSocket.GetStream();
           
            StreamReader sr = new StreamReader(readStream, Encoding.Default);
            strHTML = sr.ReadToEnd();
            String[] responseList = strHTML.Split('\n');
            if (responseList.Length >= 2)
            {
                strHTML = responseList[responseList.Length - 1];
            }
            readStream.Close();
            clientSocket.Close();

            return strHTML;
        }
    }
}

